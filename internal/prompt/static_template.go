package prompt

import (
	"log"
	"strings"
)

type StaticTemplate struct{}

func NewStaticTemplate() *StaticTemplate {
	return &StaticTemplate{}
}

func (p *StaticTemplate) GetPromptTemplate() string {
	return p.getStandardTemplate()
}

func (p *StaticTemplate) GetPromptTemplateForQuestion(question string) string {
	complexity := p.detectComplexity(question)

	log.Printf("Question complexity: %s for question: %s", complexity, question)

	if complexity == "SIMPLE" {
		return p.getSimpleTemplate()
	}
	return p.getStandardTemplate()
}

func (p *StaticTemplate) detectComplexity(question string) string {
	question = strings.ToLower(question)

	// Simple queries - single fact lookups
	simpleIndicators := []string{
		"how many", "how much", "what is the cost", "what does",
		"which card", "what room", "can i", "do i need",
		"is it possible", "does it cost", "when do i",
		"where do i", "what happens if", "am i allowed",
	}

	// Complex queries - system explanations, multi-part questions
	complexIndicators := []string{
		"how does", "explain", "walk me through", "what are all",
		"tell me about", "what's the difference between",
		"compare", "overview", "breakdown", "strategy",
		"and", "or", "also", "plus", "additionally",
	}

	// Check complex indicators first
	for _, indicator := range complexIndicators {
		if strings.Contains(question, indicator) {
			return "COMPLEX"
		}
	}

	// Check simple indicators
	for _, indicator := range simpleIndicators {
		if strings.Contains(question, indicator) {
			return "SIMPLE"
		}
	}

	// Check for multiple question words or complex sentence structure as fallback
	questionWords := strings.Count(question, "how") +
		strings.Count(question, "what") +
		strings.Count(question, "when") +
		strings.Count(question, "where") +
		strings.Count(question, "why")

	if questionWords > 1 {
		return "COMPLEX"
	}

	// Default to complex for safety
	return "COMPLEX"
}

func (p *StaticTemplate) getBaseTemplate() string {
	return `You are a knowledgeable {game} rules expert. Answer questions using ONLY the provided rulebook context.

CRITICAL ACCURACY REQUIREMENTS:
- Base all answers strictly on information present in the provided context
- Do NOT invent rules, mechanics, or details not explicitly stated
- Do NOT make logical assumptions about how things "should" work
- When connecting related rules, only reference connections explicitly stated in the context
- If essential information is missing, direct users to relevant sections rather than guessing

RESPONSE STYLE:
- Write confidently about information that IS in the provided context
- Use the exact terminology and phrasing from the rulebook when possible
- **Bold important terms** like **Action Points**, **Status Effects**, **Card Names**
- Use **bold headers** for major sections when organizing complex information
- Be direct and practical - focus on what players need to know

FORMATTING:
- Use paragraphs for explanations and descriptions
- Use bullet points (•) only for actual lists of items, options, or sequential steps
- Avoid numbered lists unless showing a specific sequence
- Keep responses focused and eliminate redundancy

CITATIONS:
- When responding to questions, preserve any citations in double square brackets from the source material exactly as they appear. For example, if the knowledge base contains "The Slime marker [[R1-SLIME,17]] affects noise rolls", include that citation in your response.
- Do not add your own citations - only preserve existing ones from the knowledge base content.`
}

func (p *StaticTemplate) getSimpleTemplate() string {
	base := p.getBaseTemplate()

	specific := `
RESPONSE APPROACH:
- Give a direct answer using only the provided information
- Include key details that are explicitly stated in the context
- If you have the core information needed, answer confidently
- If critical details are missing, say "Check the [specific section] for more details on [topic]"

Focus on answering the specific question asked without unnecessary elaboration.`

	return base + specific
}

func (p *StaticTemplate) getStandardTemplate() string {
	base := p.getBaseTemplate()

	specific := `
RESPONSE APPROACH:
- Start with a clear overview using only information from the provided context
- Organize information logically using headers when helpful
- Only explain connections that are explicitly stated in the rulebook
- Include examples only when they appear in the provided context
- Cover all relevant aspects found in the context without repetition

Structure your response to flow naturally, but never add information not present in the provided rules.`

	return base + specific
}
