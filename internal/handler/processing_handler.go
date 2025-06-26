package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

type ProcessingRequest struct {
	GameName string `json:"game_name"`
	Force    bool   `json:"force,omitempty"`
}

type ProcessingHandler struct {
	processor *knowledge.Processor
}

func NewProcessingHandler(processor *knowledge.Processor) *ProcessingHandler {
	return &ProcessingHandler{
		processor: processor,
	}
}

func (h *ProcessingHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := h.parseAndValidateRequest(request.Body)
	if err != nil {
		return utils.CreateErrorResponse(400, err.Error()), nil
	}

	result, err := h.processor.ProcessGame(ctx, req.GameName)
	if err != nil {
		return utils.CreateErrorResponse(500, err.Error()), nil
	}

	return utils.CreateSuccessResponse(result)
}

func (h *ProcessingHandler) parseAndValidateRequest(body string) (*ProcessingRequest, error) {
	var req ProcessingRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, fmt.Errorf("invalid request format: %v", err)
	}

	if req.GameName == "" {
		return nil, fmt.Errorf("game_name is required")
	}

	return &req, nil
}
