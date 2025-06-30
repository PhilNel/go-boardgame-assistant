package types

// AnswerRequest represents a request to generate an answer
// This is shared across multiple packages so it stays here
type AnswerRequest struct {
	GameName  string
	Knowledge string
	Question  string
}
