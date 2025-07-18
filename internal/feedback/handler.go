package feedback

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FeedbackRepository interface {
	SaveFeedback(ctx context.Context, feedback *FeedbackRecord) error
}

type Handler struct {
	feedbackRepo FeedbackRepository
}

func NewHandler(feedbackRepo FeedbackRepository) *Handler {
	return &Handler{
		feedbackRepo: feedbackRepo,
	}
}

func (h *Handler) SubmitFeedback(ctx context.Context, submission *FeedbackSubmission) (*FeedbackResponse, error) {
	log.Printf("Processing feedback submission for message: %s, game: %s, type: %s",
		submission.MessageID, submission.GameName, submission.FeedbackType)

	if err := h.validateSubmission(submission); err != nil {
		log.Printf("Validation failed for feedback submission: %v", err)
		return nil, err
	}

	feedback, err := h.createFeedbackRecord(submission)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback record: %w", err)
	}

	if err := h.feedbackRepo.SaveFeedback(ctx, feedback); err != nil {
		return nil, fmt.Errorf("failed to save feedback: %w", err)
	}

	log.Printf("Successfully processed feedback: %s for message: %s", feedback.FeedbackID, feedback.MessageID)

	return &FeedbackResponse{
		FeedbackID: feedback.FeedbackID,
		Message:    "Feedback submitted successfully",
	}, nil
}

func (h *Handler) validateSubmission(submission *FeedbackSubmission) error {
	if strings.TrimSpace(submission.MessageID) == "" {
		return &ValidationError{
			Code:    "INVALID_MESSAGE_ID",
			Message: "Message ID is required",
		}
	}

	if strings.TrimSpace(submission.GameName) == "" {
		return &ValidationError{
			Code:    "INVALID_GAME_NAME",
			Message: "Game name is required",
		}
	}

	if !ValidFeedbackTypes[submission.FeedbackType] {
		return &ValidationError{
			Code:    "INVALID_FEEDBACK_TYPE",
			Message: "Feedback type must be one of: positive, negative",
		}
	}

	// Validate negative feedback requirements
	if submission.FeedbackType == FeedbackTypeNegative {
		if len(submission.Issues) == 0 {
			return &ValidationError{
				Code:    "MISSING_ISSUES",
				Message: "Issues are required for negative feedback",
			}
		}

		for _, issue := range submission.Issues {
			if !ValidFeedbackIssues[issue] {
				return &ValidationError{
					Code:    "INVALID_ISSUE",
					Message: fmt.Sprintf("Invalid issue type: %s", issue),
				}
			}
		}
	}

	if len(submission.Description) > 256 {
		return &ValidationError{
			Code:    "DESCRIPTION_TOO_LONG",
			Message: "Description must be 256 characters or less",
		}
	}

	// Validate conversation context if provided
	if submission.ConversationContext != nil {
		if len(submission.ConversationContext.RecentQA) > 10 {
			return &ValidationError{
				Code:    "TOO_MANY_QA_PAIRS",
				Message: "Too many Q&A pairs in conversation context",
			}
		}

		for _, qa := range submission.ConversationContext.RecentQA {
			if len(qa.Question) > 500 {
				return &ValidationError{
					Code:    "QUESTION_TOO_LONG",
					Message: "Question in conversation context is too long",
				}
			}
			if len(qa.Answer) > 5000 {
				return &ValidationError{
					Code:    "ANSWER_TOO_LONG",
					Message: "Answer in conversation context is too long",
				}
			}
		}
	}

	return nil
}

func (h *Handler) createFeedbackRecord(submission *FeedbackSubmission) (*FeedbackRecord, error) {
	timestamp, err := time.Parse(time.RFC3339, submission.Timestamp)
	if err != nil {
		log.Printf("Failed to parse timestamp %s, using current time: %v", submission.Timestamp, err)
		timestamp = time.Now()
	}

	feedbackID := uuid.New().String()

	feedback := &FeedbackRecord{
		FeedbackID:          feedbackID,
		MessageID:           submission.MessageID,
		GameName:            submission.GameName,
		FeedbackType:        submission.FeedbackType,
		UserHash:            submission.UserHash,
		Issues:              submission.Issues,
		Description:         submission.Description,
		ConversationContext: submission.ConversationContext,
		Timestamp:           timestamp,
		CreatedAt:           time.Now().Unix(),
	}

	return feedback, nil
}
