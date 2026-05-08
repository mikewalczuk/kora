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

type Concept struct {
	ID      string // UUID string; empty for new concepts
	Title   string
	Content string
}

type GenerateResult struct {
	NewConcepts         []Concept
	UpdatedConcepts     []Concept // ID set; only title/content changes — FSRS state untouched
	ObsoletedConceptIDs []string  // UUIDs of concepts to archive; conservative — only clear absences
	Exercises           []Exercise
}

type Generator interface {
	Generate(ctx context.Context, noteTitle, noteContent string, existingConcepts []Concept) (GenerateResult, error)
}
