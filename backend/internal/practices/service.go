package practices

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/impez/kora/internal/ai"
	"github.com/impez/kora/internal/api"
	"github.com/impez/kora/internal/database"
	"github.com/impez/kora/internal/events"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNotFound = errors.New("practice not found")
	ErrConflict = errors.New("active practice already exists")
)

// DB-internal structures for the exercises JSONB column.
// correctOptionId is stored in DB but never returned to clients.
type dbOption struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}

type dbSubmission struct {
	Type             string    `json:"type"`
	SelectedOptionID uuid.UUID `json:"selectedOptionId"`
	Correct          bool      `json:"correct"`
}

type dbQuestion struct {
	ID              uuid.UUID     `json:"id"`
	Question        string        `json:"question"`
	Options         []dbOption    `json:"options"`
	CorrectOptionID uuid.UUID     `json:"correctOptionId"`
	Submission      *dbSubmission `json:"submission"`
}

type dbExercise struct {
	ID        uuid.UUID    `json:"id"`
	Type      string       `json:"type"`
	Questions []dbQuestion `json:"questions"`
}

type Service struct {
	DB  *database.Queries
	Hub *events.Hub
	AI  ai.Generator
}

type CreateInput struct {
	NoteID uuid.UUID
}

type ListInput struct {
	NoteID *uuid.UUID
	Status []api.ListPracticesParamsStatus
	Limit  int
	Page   int
}

func (s *Service) Create(ctx context.Context, in CreateInput) (uuid.UUID, error) {
	noteID := pgtype.UUID{Bytes: in.NoteID, Valid: true}

	existing, err := s.DB.GetActivePracticeForNote(ctx, noteID)
	if err == nil {
		return existing.ID.Bytes, ErrConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, err
	}

	note, err := s.DB.GetNoteForPractice(ctx, noteID)
	if err != nil {
		return uuid.UUID{}, err
	}

	aiExercises, err := s.AI.GenerateExercises(ctx, note.Title, note.Content)
	if err != nil {
		return uuid.UUID{}, err
	}

	exercises := make([]dbExercise, len(aiExercises))
	for i, ex := range aiExercises {
		exercises[i] = aiExerciseToDb(ex)
	}

	exercisesJSON, err := json.Marshal(exercises)
	if err != nil {
		return uuid.UUID{}, err
	}

	row, err := s.DB.CreatePractice(ctx, database.CreatePracticeParams{
		NoteID:    noteID,
		Status:    "in_progress",
		Exercises: exercisesJSON,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	practiceID := row.ID.Bytes

	go func() {
		time.Sleep(5 * time.Second)
		_ = s.Hub.Broadcast("practice_ready", map[string]any{
			"practiceId": uuid.UUID(practiceID).String(),
		})
	}()

	return practiceID, nil
}

func (s *Service) List(ctx context.Context, in ListInput) (api.ListPracticesResponse, error) {
	var noteID pgtype.UUID
	if in.NoteID != nil {
		noteID = pgtype.UUID{Bytes: *in.NoteID, Valid: true}
	}

	statuses := make([]string, len(in.Status))
	for i, st := range in.Status {
		statuses[i] = string(st)
	}

	offset := int32((in.Page - 1) * in.Limit)

	rows, err := s.DB.ListPractices(ctx, database.ListPracticesParams{
		NoteID:   noteID,
		Statuses: statuses,
		Off:      offset,
		Lim:      int32(in.Limit),
	})
	if err != nil {
		return api.ListPracticesResponse{}, err
	}

	count, err := s.DB.CountPractices(ctx, database.CountPracticesParams{
		NoteID:   noteID,
		Statuses: statuses,
	})
	if err != nil {
		return api.ListPracticesResponse{}, err
	}

	practices := make([]api.Practice, 0, len(rows))
	for _, row := range rows {
		p, err := rowToPractice(row)
		if err != nil {
			continue
		}
		practices = append(practices, p)
	}

	return api.ListPracticesResponse{Items: practices, Total: int(count)}, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (api.Practice, error) {
	row, err := s.DB.GetPractice(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return api.Practice{}, ErrNotFound
	}
	if err != nil {
		return api.Practice{}, err
	}
	return rowToPractice(row)
}

func (s *Service) Submit(ctx context.Context, practiceID, exerciseID uuid.UUID, body *api.SubmitExerciseJSONRequestBody) (api.SubmitExerciseResponse, error) {
	row, err := s.DB.GetPractice(ctx, pgtype.UUID{Bytes: practiceID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return api.SubmitExerciseResponse{}, ErrNotFound
	}
	if err != nil {
		return api.SubmitExerciseResponse{}, err
	}

	var exercises []dbExercise
	if err := json.Unmarshal(row.Exercises, &exercises); err != nil {
		return api.SubmitExerciseResponse{}, err
	}

	sub, err := body.AsMultiQuizSubmission()
	if err != nil {
		return api.SubmitExerciseResponse{}, err
	}

	var result api.MultiQuizExerciseResult
	found := false

	for ei := range exercises {
		if exercises[ei].ID != exerciseID {
			continue
		}
		for qi := range exercises[ei].Questions {
			if exercises[ei].Questions[qi].ID != sub.QuestionId {
				continue
			}
			q := &exercises[ei].Questions[qi]
			if q.Submission != nil {
				return api.SubmitExerciseResponse{}, ErrConflict
			}
			correct := q.CorrectOptionID == sub.SelectedOptionId
			q.Submission = &dbSubmission{
				Type:             "multi-quiz",
				SelectedOptionID: sub.SelectedOptionId,
				Correct:          correct,
			}
			result = api.MultiQuizExerciseResult{
				Type:             api.MultiQuizExerciseResultTypeMultiQuiz,
				QuestionId:       sub.QuestionId,
				SelectedOptionId: sub.SelectedOptionId,
				Correct:          correct,
			}
			found = true
			break
		}
		break
	}

	if !found {
		return api.SubmitExerciseResponse{}, ErrNotFound
	}

	exercisesJSON, err := json.Marshal(exercises)
	if err != nil {
		return api.SubmitExerciseResponse{}, err
	}

	_, err = s.DB.UpdatePracticeExercises(ctx, database.UpdatePracticeExercisesParams{
		ID:        pgtype.UUID{Bytes: practiceID, Valid: true},
		Exercises: exercisesJSON,
	})
	if err != nil {
		return api.SubmitExerciseResponse{}, err
	}

	if allAnswered(exercises) {
		_, _ = s.DB.CompletePractice(ctx, pgtype.UUID{Bytes: practiceID, Valid: true})
	}

	var response api.SubmitExerciseResponse
	err = response.FromMultiQuizExerciseResult(result)
	return response, err
}

func rowToPractice(row database.Practice) (api.Practice, error) {
	var exercises []dbExercise
	if err := json.Unmarshal(row.Exercises, &exercises); err != nil {
		return api.Practice{}, err
	}

	apiExercises := make([]api.Exercise, 0, len(exercises))
	for _, ex := range exercises {
		apiEx, err := dbExerciseToAPI(ex)
		if err != nil {
			return api.Practice{}, err
		}
		apiExercises = append(apiExercises, apiEx)
	}

	var completedAt *time.Time
	if row.CompletedAt.Valid {
		t := row.CompletedAt.Time
		completedAt = &t
	}

	return api.Practice{
		Id:          row.ID.Bytes,
		NoteId:      row.NoteID.Bytes,
		Status:      api.PracticeStatus(row.Status),
		Exercises:   apiExercises,
		CreatedAt:   row.CreatedAt.Time,
		CompletedAt: completedAt,
	}, nil
}

// dbExerciseToAPI converts a DB exercise to the API type, stripping
// correctOptionId but preserving submission state.
func dbExerciseToAPI(ex dbExercise) (api.Exercise, error) {
	type apiOption struct {
		Id   uuid.UUID `json:"id"`
		Text string    `json:"text"`
	}
	type apiQuestion struct {
		Id         uuid.UUID     `json:"id"`
		Question   string        `json:"question"`
		Options    []apiOption   `json:"options"`
		Submission *dbSubmission `json:"submission"`
	}
	type apiExerciseShape struct {
		Id        uuid.UUID     `json:"id"`
		Type      string        `json:"type"`
		Questions []apiQuestion `json:"questions"`
	}

	questions := make([]apiQuestion, len(ex.Questions))
	for i, q := range ex.Questions {
		opts := make([]apiOption, len(q.Options))
		for j, o := range q.Options {
			opts[j] = apiOption{Id: o.ID, Text: o.Text}
		}
		questions[i] = apiQuestion{
			Id:         q.ID,
			Question:   q.Question,
			Options:    opts,
			Submission: q.Submission,
		}
	}

	raw, err := json.Marshal(apiExerciseShape{Id: ex.ID, Type: ex.Type, Questions: questions})
	if err != nil {
		return api.Exercise{}, err
	}
	var e api.Exercise
	err = json.Unmarshal(raw, &e)
	return e, err
}

func allAnswered(exercises []dbExercise) bool {
	for _, ex := range exercises {
		for _, q := range ex.Questions {
			if q.Submission == nil {
				return false
			}
		}
	}
	return true
}

func aiExerciseToDb(ex ai.Exercise) dbExercise {
	questions := make([]dbQuestion, len(ex.Questions))
	for i, q := range ex.Questions {
		opts := make([]dbOption, len(q.Options))
		for j, o := range q.Options {
			opts[j] = dbOption{ID: uuid.New(), Text: o.Text}
		}
		questions[i] = dbQuestion{
			ID:              uuid.New(),
			Question:        q.Text,
			Options:         opts,
			CorrectOptionID: opts[q.CorrectOptionIdx].ID,
		}
	}
	return dbExercise{
		ID:        uuid.New(),
		Type:      "multi-quiz",
		Questions: questions,
	}
}
