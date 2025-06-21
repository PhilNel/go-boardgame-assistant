package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PhilNel/go-boardgame-assistant/internal/logger"
	"github.com/PhilNel/go-boardgame-assistant/internal/types"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

type Request struct {
	GameName string `json:"gameName"`
	Question string `json:"question"`
}

type Response struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

type KnowledgeProvider interface {
	GetKnowledge(ctx context.Context, gameName string) (string, error)
}

type AnswerProvider interface {
	GenerateAnswer(ctx context.Context, request *types.AnswerRequest) (string, error)
}

type QuestionHandler struct {
	knowledgeProvider KnowledgeProvider
	answerProvider    AnswerProvider
}

func NewQuestionHandler(knowledgeProvider KnowledgeProvider, answerProvider AnswerProvider) *QuestionHandler {
	return &QuestionHandler{
		knowledgeProvider: knowledgeProvider,
		answerProvider:    answerProvider,
	}
}

func (h *QuestionHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := h.parseAndValidateRequest(request.Body)
	if err != nil {
		return utils.CreateErrorResponse(400, err.Error()), nil
	}

	answer, err := h.processQuestion(ctx, req)
	if err != nil {
		return utils.CreateErrorResponse(500, err.Error()), nil
	}

	response := Response{Answer: answer}
	return utils.CreateSuccessResponse(response)
}

func (h *QuestionHandler) parseAndValidateRequest(body string) (*Request, error) {
	var req Request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, fmt.Errorf("invalid request format: %v", err)
	}

	if err := h.validateRequest(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *QuestionHandler) validateRequest(req *Request) error {
	if req.GameName == "" {
		return fmt.Errorf("gameName is required")
	}
	if req.Question == "" {
		return fmt.Errorf("question is required")
	}
	return nil
}

func (h *QuestionHandler) processQuestion(ctx context.Context, req *Request) (string, error) {
	logger.LogIncomingRequest(req.GameName, req.Question)

	knowledge, err := h.knowledgeProvider.GetKnowledge(ctx, req.GameName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve game knowledge: %w", err)
	}

	answerRequest := &types.AnswerRequest{
		GameName:  req.GameName,
		Knowledge: knowledge,
		Question:  req.Question,
	}

	answer, err := h.answerProvider.GenerateAnswer(ctx, answerRequest)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	logger.LogSuccessfulQAPair(req.GameName, req.Question, answer)

	return answer, nil
}
