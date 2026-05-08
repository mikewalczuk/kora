package ai

import "context"

// OpenAIGenerator calls a real AI provider to generate concepts and exercises from note content.
type OpenAIGenerator struct {
	APIKey string
}

func (g *OpenAIGenerator) Generate(_ context.Context, _, _ string, _ []Concept) (GenerateResult, error) {
	panic("not implemented")
}
