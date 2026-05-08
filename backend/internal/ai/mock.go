package ai

import "context"

type MockGenerator struct{}

func (m *MockGenerator) GenerateExercises(_ context.Context, _, _ string) ([]Exercise, error) {
	return []Exercise{
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
	}, nil
}
