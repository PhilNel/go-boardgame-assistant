package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PhilNel/go-boardgame-assistant/internal/feedback"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

type FeedbackHandler struct {
	feedbackHandler *feedback.Handler
}

func NewFeedbackHandler(feedbackHandler *feedback.Handler) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackHandler: feedbackHandler,
	}
}

func (h *FeedbackHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return utils.CreateErrorResponse(405, "Method not allowed"), nil
	}

	submission, err := h.parseAndValidateRequest(request.Body)
	if err != nil {
		return utils.CreateErrorResponse(400, err.Error()), nil
	}

	response, err := h.feedbackHandler.SubmitFeedback(ctx, submission)
	if err != nil {
		if validationErr, ok := err.(*feedback.ValidationError); ok {
			return utils.CreateErrorResponse(400, validationErr.Message), nil
		}
		return utils.CreateErrorResponse(500, err.Error()), nil
	}

	return utils.CreateSuccessResponse(response)
}

func (h *FeedbackHandler) parseAndValidateRequest(body string) (*feedback.FeedbackSubmission, error) {
	var submission feedback.FeedbackSubmission
	if err := json.Unmarshal([]byte(body), &submission); err != nil {
		return nil, fmt.Errorf("invalid request format: %v", err)
	}

	return &submission, nil
}
