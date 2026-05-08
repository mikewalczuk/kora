package practices

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/impez/kora/internal/api"
)

var ErrNotFound = errors.New("practice not found")

type Service struct {
	mu    sync.RWMutex
	store map[uuid.UUID]api.Practice
}

type CreateInput struct {
	NoteID uuid.UUID
}

func (s *Service) Create(ctx context.Context, in CreateInput) (api.Practice, error) {
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
		Id:    uuid.New(),
		NoteId: in.NoteID,
		Status: api.InProgress,
		Exercises: []api.Exercise{
			toExercise(api.MultiQuizExercise{Id: uuid.New(), Type: api.MultiQuiz, Questions: []api.MultiQuizQuestion{question}}),
			toExercise(api.MultiQuizExercise{Id: uuid.New(), Type: api.MultiQuiz, Questions: []api.MultiQuizQuestion{question2}}),
		},
		CreatedAt:   now,
		CompletedAt: nil,
	}

	s.mu.Lock()
	if s.store == nil {
		s.store = make(map[uuid.UUID]api.Practice)
	}
	s.store[practice.Id] = practice
	s.mu.Unlock()

	return practice, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (api.Practice, error) {
	s.mu.RLock()
	practice, ok := s.store[id]
	s.mu.RUnlock()
	if !ok {
		return api.Practice{}, ErrNotFound
	}
	return practice, nil
}
