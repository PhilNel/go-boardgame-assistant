package answer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/types"
)

type TemplateProvider interface {
	GetPromptTemplate() string
	GetPromptTemplateForQuestion(question string) string
}

type BedrockProvider struct {
	bedrockClient    aws.BedrockClient
	templateProvider TemplateProvider
	config           *config.Bedrock
}

func NewBedrockProvider(bedrockClient aws.BedrockClient, templateProvider TemplateProvider, config *config.Bedrock) *BedrockProvider {
	return &BedrockProvider{
		bedrockClient:    bedrockClient,
		templateProvider: templateProvider,
		config:           config,
	}
}

func (b *BedrockProvider) GenerateAnswer(ctx context.Context, request *types.AnswerRequest) (string, error) {
	template := b.templateProvider.GetPromptTemplateForQuestion(request.Question)
	systemPrompt := strings.ReplaceAll(template, "{game}", request.GameName)
	userContent := fmt.Sprintf("%s\n\nGame Context:\n%s\n\nQuestion: %s", systemPrompt, request.Knowledge, request.Question)

	log.Printf("Using Bedrock model ID: %s for game: %s", b.bedrockClient.GetModelID(), request.GameName)
	log.Printf("Request context: Knowledge length=%d, Question length=%d", len(request.Knowledge), len(request.Question))

	return b.generateAnswerWithClaude(ctx, userContent)
}

func (b *BedrockProvider) generateAnswerWithClaude(ctx context.Context, userContent string) (string, error) {
	bedrockRequest := &aws.BedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []aws.BedrockMessage{
			{
				Role:    "user",
				Content: userContent,
			},
		},
		MaxTokens:   b.config.AnswerMaxTokens,
		Temperature: b.config.AnswerTemperature,
		TopP:        b.config.AnswerTopP,
	}

	response, err := b.bedrockClient.InvokeModel(ctx, bedrockRequest)
	if err != nil {
		log.Printf("ERROR: Bedrock InvokeModel failed: %v", err)
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	log.Printf("Bedrock InvokeModel succeeded, extracting response...")
	answer, err := b.extractTextFromResponse(response)
	if err != nil {
		log.Printf("ERROR: Failed to extract text from response: %v", err)
		return "", err
	}

	log.Printf("Successfully extracted answer with length: %d", len(answer))
	return answer, nil
}

func (b *BedrockProvider) extractTextFromResponse(response *aws.BedrockResponse) (string, error) {
	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	var answer string
	for _, content := range response.Content {
		if content.Type == "text" {
			answer += content.Text
		}
	}

	if answer == "" {
		return "", fmt.Errorf("no text content found in response")
	}

	return answer, nil
}
