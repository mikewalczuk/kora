package ai

import "context"

// OpenAIGenerator calls a real AI provider to generate exercises from note content.
type OpenAIGenerator struct {
	APIKey string
}

func (g *OpenAIGenerator) GenerateExercises(_ context.Context, _, _ string) ([]Exercise, error) {
	panic("not implemented")
}
