package ai

import "context"

type MockGenerator struct{}

func (m *MockGenerator) Generate(_ context.Context, _, _ string, _ []Concept) (GenerateResult, error) {
	return GenerateResult{
		NewConcepts: []Concept{
			{Title: "Main Concept", Content: "The primary idea explored in this note."},
			{Title: "Key Takeaway", Content: "The most important conclusion drawn from the material."},
		},
		UpdatedConcepts:     nil,
		ObsoletedConceptIDs: nil,
		Exercises: []Exercise{
			{
				Questions: []Question{
					{
						Text:             "What is the main concept discussed in this note?",
						Options:          []Option{{Text: "Option A"}, {Text: "Option B"}, {Text: "Option C"}},
						CorrectOptionIdx: 0,
					},
				},
			},
			{
				Questions: []Question{
					{
						Text:             "Which of the following best summarizes the key takeaway?",
						Options:          []Option{{Text: "Option A"}, {Text: "Option B"}, {Text: "Option C"}},
						CorrectOptionIdx: 0,
					},
				},
			},
		},
	}, nil
}
