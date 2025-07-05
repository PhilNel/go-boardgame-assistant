package prompt

import (
	"fmt"
	"log"
	"strings"
)

type StaticTemplate struct{}

func NewStaticTemplate() *StaticTemplate {
	return &StaticTemplate{}
}

func (p *StaticTemplate) GetPromptTemplate() string {
	return p.getBroadQuestionTemplate()
}

func (p *StaticTemplate) GetPromptTemplateForQuestion(question string) string {
	questionType := p.detectQuestionType(question)

	log.Printf("Question type detected: %s for question: %s", questionType, question)

	if questionType == "NARROW" {
		return p.getNarrowQuestionTemplate()
	}
	return p.getBroadQuestionTemplate()
}

func (p *StaticTemplate) detectQuestionType(question string) string {
	question = strings.ToLower(question)

	// Broad question patterns - system explanations and overviews
	broadPatterns := []string{
		"how does", "how do", "explain", "what is", "tell me about",
		"work", "system", "mechanic", "overview", "breakdown",
		"walk me through", "give me", "describe", "what are all",
	}

	// Narrow question patterns - specific queries and yes/no questions
	narrowPatterns := []string{
		"how are", "how is", "what happens when", "can i", "can you",
		"do i", "does", "is it possible", "what causes", "when does",
		"where do", "which", "what deck", "how many", "what room",
		"should i", "must i", "am i allowed", "is there a way",
	}

	// Check broad patterns first to catch "how does" before "does"
	for _, pattern := range broadPatterns {
		if strings.Contains(question, pattern) {
			return "BROAD"
		}
	}

	for _, pattern := range narrowPatterns {
		if strings.Contains(question, pattern) {
			return "NARROW"
		}
	}

	return "BROAD" // Default to comprehensive for unclear cases
}

func (p *StaticTemplate) getCommonTemplate() string {
	return `You are an expert on {game} board game rules. %s

CORE PRINCIPLES:
- Answer ONLY using information explicitly stated in the provided context
- Do NOT add details, assumptions, or logical inferences not directly stated in the rules
- If specific mechanics aren't explained in the context, simply state what IS covered
- Stick strictly to the exact wording and information provided

RESPONSE FORMATTING:
- **Bold section headers** like "**Combat Resolution**", "**Movement Rules**", "**Victory Conditions**"
- Use bullet points (•) only for actual lists of items, steps, or key points
- Write explanatory content as natural paragraphs, not bullet points
- **Bold important game terms** like **Action Points**, **Status Effects**, **Card Types**
- Avoid numbered lists - use bullet points sparingly and only when listing discrete items

TONE AND STYLE:
- %s
- Use clear, natural language transitions
- Be conversational but authoritative
- Only include examples that are explicitly mentioned in the provided context
- Never invent steps, procedures, or details not stated in the rules`
}

func (p *StaticTemplate) getNarrowQuestionTemplate() string {
	common := p.getCommonTemplate()

	specificApproach := `FOCUSED RESPONSE STRATEGY:
• Lead with a direct answer using only the provided information
• List only the specific details explicitly stated in the context
• Do NOT add implied steps, assumed procedures, or invented details
• If the context doesn't fully answer the question, state what IS covered without padding

STRUCTURE:
- **Direct Answer**: Clear, immediate response to the question
- **Key Details**: Specific rules or mechanics (use paragraphs for explanations, bullets only for lists)
- **Context** (if needed): Brief related information
- **Practical Note**: How this applies in gameplay`

	return fmt.Sprintf(common,
		"Answer the specific question asked using the provided knowledge base.",
		"Start with a direct, clear answer to the question") + "\n\n" + specificApproach
}

func (p *StaticTemplate) getBroadQuestionTemplate() string {
	common := p.getCommonTemplate()

	specificApproach := `COMPREHENSIVE RESPONSE STRATEGY:
• Begin with an overview using only information from the provided context
• Break down only the components explicitly mentioned in the rules
• Explain relationships only when they are directly stated
• Include only examples and cases explicitly mentioned in the context
• Only mention connections that are directly stated in the rules

STRUCTURE:
- **Overview**: Brief introduction using only provided information (write as paragraphs)
- **Core Components**: Only elements explicitly mentioned in the context (use bullets for lists, paragraphs for explanations)
- **Key Rules**: Only mechanics and restrictions directly stated (write as paragraphs unless listing multiple items)
- **Examples**: Only scenarios explicitly mentioned in the rules
- **Related Systems**: Only connections explicitly stated in the context`

	return fmt.Sprintf(common,
		"Provide a comprehensive explanation using the provided knowledge base.",
		"Begin with a clear overview of the system or mechanic") + "\n\n" + specificApproach
}
