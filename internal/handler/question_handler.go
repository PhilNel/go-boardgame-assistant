package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
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
	GetKnowledge(ctx context.Context, gameName string, query string) (string, error)
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

	knowledgeContent, err := h.knowledgeProvider.GetKnowledge(ctx, req.GameName, req.Question)
	if err != nil {
		// Check if this is a "no relevant knowledge" error
		var noKnowledgeErr *knowledge.NoRelevantKnowledgeError
		if errors.As(err, &noKnowledgeErr) {
			// Return a helpful message instead of an error
			return "I don't have any specific information about that topic in my knowledge base for " + req.GameName +
				". This might be something we haven't covered yet, or your question might need to be more specific. " +
				"Feel free to try rephrasing your question or asking about a different aspect of the game!", nil
		}
		return "", fmt.Errorf("failed to retrieve game knowledge: %w", err)
	}

	answerRequest := &types.AnswerRequest{
		GameName:  req.GameName,
		Knowledge: knowledgeContent,
		Question:  req.Question,
	}

	answer, err := h.answerProvider.GenerateAnswer(ctx, answerRequest)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	logger.LogSuccessfulQAPair(req.GameName, req.Question, answer)

	return answer, nil
}
