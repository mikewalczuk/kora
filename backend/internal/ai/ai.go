package ai

import "context"

type Option struct {
	Text string
}

type Question struct {
	Text             string
	Options          []Option
	CorrectOptionIdx int
}

type Exercise struct {
	Questions []Question
}

type Generator interface {
	GenerateExercises(ctx context.Context, noteTitle, noteContent string) ([]Exercise, error)
}
