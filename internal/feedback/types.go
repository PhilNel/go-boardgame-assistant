package feedback

import (
	"time"
)

type FeedbackType string

const (
	FeedbackTypePositive FeedbackType = "positive"
	FeedbackTypeNegative FeedbackType = "negative"
)

type FeedbackIssue string

const (
	FeedbackIssueIncorrectInfo FeedbackIssue = "incorrect_info"
	FeedbackIssueMissingInfo   FeedbackIssue = "missing_info"
	FeedbackIssueUnclear       FeedbackIssue = "unclear"
	FeedbackIssueWrongGame     FeedbackIssue = "wrong_game"
	FeedbackIssueOther         FeedbackIssue = "other"
)

type QAPair struct {
	Question  string `json:"question" dynamodbav:"question"`
	Answer    string `json:"answer" dynamodbav:"answer"`
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
}

type ConversationContext struct {
	RecentQA []QAPair `json:"recent_qa" dynamodbav:"recent_qa"`
}

// FeedbackRecord represents a feedback submission stored in the database
type FeedbackRecord struct {
	FeedbackID          string               `json:"feedback_id" dynamodbav:"feedback_id"`
	MessageID           string               `json:"message_id" dynamodbav:"message_id"`
	GameName            string               `json:"game_name" dynamodbav:"game_name"`
	FeedbackType        FeedbackType         `json:"feedback_type" dynamodbav:"feedback_type"`
	UserHash            string               `json:"user_hash,omitempty" dynamodbav:"user_hash"`
	Issues              []FeedbackIssue      `json:"issues,omitempty" dynamodbav:"issues"`
	Description         string               `json:"description,omitempty" dynamodbav:"description"`
	ConversationContext *ConversationContext `json:"conversation_context,omitempty" dynamodbav:"conversation_context"`
	Timestamp           time.Time            `json:"timestamp" dynamodbav:"timestamp"`
	CreatedAt           int64                `json:"created_at" dynamodbav:"created_at"`
}

// FeedbackSubmission represents the incoming feedback submission request
type FeedbackSubmission struct {
	MessageID           string               `json:"message_id"`
	GameName            string               `json:"game_name"`
	FeedbackType        FeedbackType         `json:"feedback_type"`
	UserHash            string               `json:"user_hash,omitempty"`
	Issues              []FeedbackIssue      `json:"issues,omitempty"`
	Description         string               `json:"description,omitempty"`
	ConversationContext *ConversationContext `json:"conversation_context,omitempty"`
	Timestamp           string               `json:"timestamp"`
}

// FeedbackResponse represents the response after submitting feedback
type FeedbackResponse struct {
	FeedbackID string `json:"feedback_id"`
	Message    string `json:"message"`
}

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

var ValidFeedbackTypes = map[FeedbackType]bool{
	FeedbackTypePositive: true,
	FeedbackTypeNegative: true,
}

var ValidFeedbackIssues = map[FeedbackIssue]bool{
	FeedbackIssueIncorrectInfo: true,
	FeedbackIssueMissingInfo:   true,
	FeedbackIssueUnclear:       true,
	FeedbackIssueWrongGame:     true,
	FeedbackIssueOther:         true,
}
