package practices

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/impez/kora/internal/api"
	"github.com/impez/kora/internal/events"
)

var (
	ErrNotFound = errors.New("practice not found")
	ErrConflict = errors.New("active practice already exists")
)

type Service struct {
	Hub   *events.Hub
	mu    sync.RWMutex
	store map[uuid.UUID]api.Practice
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

func (s *Service) List(_ context.Context, in ListInput) (api.ListPracticesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []api.Practice
	for _, p := range s.store {
		if in.NoteID != nil && p.NoteId != *in.NoteID {
			continue
		}
		if len(in.Status) > 0 && !statusMatch(in.Status, p.Status) {
			continue
		}
		matches = append(matches, p)
	}

	total := len(matches)
	start := (in.Page - 1) * in.Limit
	if start >= total {
		matches = nil
	} else {
		end := start + in.Limit
		if end > total {
			end = total
		}
		matches = matches[start:end]
	}
	if matches == nil {
		matches = []api.Practice{}
	}

	return api.ListPracticesResponse{Items: matches, Total: total}, nil
}

func statusMatch(filter []api.ListPracticesParamsStatus, s api.PracticeStatus) bool {
	for _, f := range filter {
		if string(f) == string(s) {
			return true
		}
	}
	return false
}

func (s *Service) Create(ctx context.Context, in CreateInput) (uuid.UUID, error) {
	s.mu.Lock()
	if s.store == nil {
		s.store = make(map[uuid.UUID]api.Practice)
	}
	for _, p := range s.store {
		if p.NoteId == in.NoteID &&
			(p.Status == api.PracticeStatusInProgress || p.Status == api.PracticeStatusPending) {
			s.mu.Unlock()
			return p.Id, ErrConflict
		}
	}

	question := api.MultiQuizQuestion{
		Id:       uuid.New(),
		Question: "What is the main concept discussed in this note?",
		Options: []struct {
			Id   uuid.UUID `json:"id"`
			Text string    `json:"text"`
		}{
			{Id: uuid.New(), Text: "Option A"},
			{Id: uuid.New(), Text: "Option B"},
			{Id: uuid.New(), Text: "Option C"},
		},
	}
	question2 := api.MultiQuizQuestion{
		Id:       uuid.New(),
		Question: "Which of the following best summarizes the key takeaway?",
		Options: []struct {
			Id   uuid.UUID `json:"id"`
			Text string    `json:"text"`
		}{
			{Id: uuid.New(), Text: "Option A"},
			{Id: uuid.New(), Text: "Option B"},
			{Id: uuid.New(), Text: "Option C"},
		},
	}

	toExercise := func(mq api.MultiQuizExercise) api.Exercise {
		raw, _ := json.Marshal(mq)
		var e api.Exercise
		_ = json.Unmarshal(raw, &e)
		return e
	}

	now := time.Now().UTC()
	practice := api.Practice{
		Id:     uuid.New(),
		NoteId: in.NoteID,
		Status: api.PracticeStatusInProgress,
		Exercises: []api.Exercise{
			toExercise(api.MultiQuizExercise{Id: uuid.New(), Type: api.MultiQuiz, Questions: []api.MultiQuizQuestion{question}}),
			toExercise(api.MultiQuizExercise{Id: uuid.New(), Type: api.MultiQuiz, Questions: []api.MultiQuizQuestion{question2}}),
		},
		CreatedAt:   now,
		CompletedAt: nil,
	}
	s.store[practice.Id] = practice
	s.mu.Unlock()

	go func() {
		time.Sleep(5 * time.Second)
		_ = s.Hub.Broadcast("practice_ready", map[string]any{
			"practiceId": practice.Id,
		})
	}()

	return practice.Id, nil
}

func (s *Service) Get(_ context.Context, id uuid.UUID) (api.Practice, error) {
	s.mu.RLock()
	practice, ok := s.store[id]
	s.mu.RUnlock()
	if !ok {
		return api.Practice{}, ErrNotFound
	}
	return practice, nil
}
