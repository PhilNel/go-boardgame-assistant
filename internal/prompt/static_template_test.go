package prompt

import (
	"strings"
	"testing"
)

func TestDetectQuestionType(t *testing.T) {
	template := NewStaticTemplate()

	testCases := []struct {
		question     string
		expectedType string
		description  string
	}{
		// Narrow question tests
		{"How are players damaged by fire?", "NARROW", "how are pattern"},
		{"How is movement calculated?", "NARROW", "how is pattern"},
		{"What happens when a player dies?", "NARROW", "what happens when pattern"},
		{"Can I move through walls?", "NARROW", "can i pattern"},
		{"Can you attack multiple enemies?", "NARROW", "can you pattern"},
		{"Do I need to roll dice?", "NARROW", "do i pattern"},
		{"Does combat use initiative?", "NARROW", "does pattern"},
		{"Is it possible to heal?", "NARROW", "is it possible pattern"},
		{"What causes fire damage?", "NARROW", "what causes pattern"},
		{"When does the game end?", "NARROW", "when does pattern"},
		{"Where do I place tokens?", "NARROW", "where do pattern"},
		{"Which cards can I play?", "NARROW", "which pattern"},
		{"What deck do I draw from?", "NARROW", "what deck pattern"},
		{"How many actions do I get?", "NARROW", "how many pattern"},
		{"What room am I in?", "NARROW", "what room pattern"},

		// Broad question tests
		{"How does combat work?", "BROAD", "how does pattern"},
		{"How do players move?", "BROAD", "how do pattern"},
		{"Explain the magic system", "BROAD", "explain pattern"},
		{"What is the victory condition?", "BROAD", "what is pattern"},
		{"Tell me about character creation", "BROAD", "tell me about pattern"},
		{"How does the trading system work?", "BROAD", "work pattern"},
		{"What is the combat system?", "BROAD", "system pattern"},
		{"Explain the movement mechanic", "BROAD", "mechanic pattern"},
		{"Give me an overview of the game", "BROAD", "overview pattern"},
		{"Provide a breakdown of turns", "BROAD", "breakdown pattern"},

		// Edge cases and unclear questions (should default to BROAD)
		{"What should I do next?", "BROAD", "unclear question defaults to broad"},
		{"Help me understand this", "BROAD", "vague question defaults to broad"},
		{"I'm confused", "BROAD", "very vague defaults to broad"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := template.detectQuestionType(tc.question)
			if result != tc.expectedType {
				t.Errorf("Question: %q\nExpected: %s\nGot: %s\nDescription: %s",
					tc.question, tc.expectedType, result, tc.description)
			}
		})
	}
}

func TestGetPromptTemplateForQuestion(t *testing.T) {
	template := NewStaticTemplate()

	// Test narrow question returns narrow template
	narrowQuestion := "How are players damaged?"
	narrowTemplate := template.GetPromptTemplateForQuestion(narrowQuestion)
	if !containsNarrowTemplateMarkers(narrowTemplate) {
		t.Error("Narrow question should return narrow template")
	}

	// Test broad question returns broad template
	broadQuestion := "How does combat work?"
	broadTemplate := template.GetPromptTemplateForQuestion(broadQuestion)
	if !containsBroadTemplateMarkers(broadTemplate) {
		t.Error("Broad question should return broad template")
	}

	// Verify templates are different
	if narrowTemplate == broadTemplate {
		t.Error("Narrow and broad templates should be different")
	}
}

func containsNarrowTemplateMarkers(template string) bool {
	// Check for markers that indicate this is the narrow template
	return strings.Contains(template, "STRUCTURE FOR FOCUSED QUESTIONS") &&
		strings.Contains(template, "Don't provide comprehensive explanations for narrow questions")
}

func containsBroadTemplateMarkers(template string) bool {
	// Check for markers that indicate this is the broad template
	return strings.Contains(template, "STRUCTURE FOR COMPREHENSIVE QUESTIONS") &&
		strings.Contains(template, "Provide a thorough explanation of the system or mechanic")
}
