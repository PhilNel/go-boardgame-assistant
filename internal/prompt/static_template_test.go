package prompt

import (
	"strings"
	"testing"
)

func TestDetectComplexity(t *testing.T) {
	template := NewStaticTemplate()

	testCases := []struct {
		question     string
		expectedType string
		description  string
	}{
		// Simple question tests
		{"How many actions do I get?", "SIMPLE", "how many pattern"},
		{"How much does it cost?", "SIMPLE", "how much pattern"},
		{"What is the cost of this card?", "SIMPLE", "what is the cost pattern"},
		{"What does this ability do?", "SIMPLE", "what does pattern"},
		{"Which card should I play?", "SIMPLE", "which card pattern"},
		{"What room am I in?", "SIMPLE", "what room pattern"},
		{"Can I move through walls?", "SIMPLE", "can i pattern"},
		{"Do I need to roll dice?", "SIMPLE", "do i need pattern"},
		{"Is it possible to heal?", "SIMPLE", "is it possible pattern"},
		{"Does it cost an action?", "SIMPLE", "does it cost pattern"},
		{"When do I draw cards?", "SIMPLE", "when do i pattern"},
		{"Where do I place tokens?", "SIMPLE", "where do i pattern"},
		{"What happens if I fail?", "SIMPLE", "what happens if pattern"},
		{"Am I allowed to attack?", "SIMPLE", "am i allowed pattern"},

		// Complex question tests
		{"How does combat work?", "COMPLEX", "how does pattern"},
		{"How do players move?", "COMPLEX", "how does pattern"},
		{"Explain the magic system", "COMPLEX", "explain pattern"},
		{"Walk me through character creation", "COMPLEX", "walk me through pattern"},
		{"Tell me about the trading system", "COMPLEX", "tell me about pattern"},
		{"What are all the victory conditions?", "COMPLEX", "what are all pattern"},
		{"What's the difference between melee and ranged?", "COMPLEX", "what's the difference pattern"},
		{"Compare these two abilities", "COMPLEX", "compare pattern"},
		{"Give me an overview of the game", "COMPLEX", "overview pattern"},
		{"Provide a breakdown of turns", "COMPLEX", "breakdown pattern"},
		{"What is the strategy for winning?", "COMPLEX", "strategy pattern"},

		// Multiple question words (should be COMPLEX)
		{"How many cards can I draw and when?", "COMPLEX", "multiple question words"},
		{"What happens when I attack and how much damage?", "COMPLEX", "multiple question words"},

		// Compound questions with conjunctions (should be COMPLEX)
		{"Can I move and attack in the same turn?", "COMPLEX", "contains 'and'"},
		{"Should I heal or attack this turn?", "COMPLEX", "contains 'or'"},
		{"Can I also use this ability?", "COMPLEX", "contains 'also'"},

		// Edge cases (should default to COMPLEX)
		{"What should I do next?", "COMPLEX", "defaults to complex"},
		{"Help me understand this", "COMPLEX", "vague question defaults to complex"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := template.detectComplexity(tc.question)
			if result != tc.expectedType {
				t.Errorf("Question: %q\nExpected: %s\nGot: %s\nDescription: %s",
					tc.question, tc.expectedType, result, tc.description)
			}
		})
	}
}

func TestGetPromptTemplateForQuestion(t *testing.T) {
	template := NewStaticTemplate()

	// Test simple question returns simple template
	simpleQuestion := "How many actions do I get?"
	simpleTemplate := template.GetPromptTemplateForQuestion(simpleQuestion)
	if !containsSimpleTemplateMarkers(simpleTemplate) {
		t.Error("Simple question should return simple template")
	}

	// Test complex question returns standard template
	complexQuestion := "How does combat work?"
	complexTemplate := template.GetPromptTemplateForQuestion(complexQuestion)
	if !containsStandardTemplateMarkers(complexTemplate) {
		t.Error("Complex question should return standard template")
	}

	// Verify templates are different
	if simpleTemplate == complexTemplate {
		t.Error("Simple and standard templates should be different")
	}
}

func containsSimpleTemplateMarkers(template string) bool {
	// Check for markers that indicate this is the simple template
	return strings.Contains(template, "Give a direct answer using only the provided information") &&
		strings.Contains(template, "Focus on answering the specific question asked without unnecessary elaboration")
}

func containsStandardTemplateMarkers(template string) bool {
	// Check for markers that indicate this is the standard template
	return strings.Contains(template, "Start with a clear overview using only information from the provided context") &&
		strings.Contains(template, "Structure your response to flow naturally")
}
