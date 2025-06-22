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
}

type BedrockProvider struct {
	bedrockClient    *aws.BedrockClient
	templateProvider TemplateProvider
}

func NewBedrockProvider(config *config.Bedrock, templateProvider TemplateProvider) (*BedrockProvider, error) {
	bedrockClient, err := aws.NewBedrockClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Bedrock client: %w", err)
	}

	return &BedrockProvider{
		bedrockClient:    bedrockClient,
		templateProvider: templateProvider,
	}, nil
}

func (b *BedrockProvider) GenerateAnswer(ctx context.Context, request *types.AnswerRequest) (string, error) {
	template := b.templateProvider.GetPromptTemplate()
	systemPrompt := strings.ReplaceAll(template, "{game}", request.GameName)
	userContent := fmt.Sprintf("%s\n\nGame Context:\n%s\n\nQuestion: %s", systemPrompt, request.Knowledge, request.Question)

	log.Printf("Using Bedrock model ID: %s for game: %s", b.bedrockClient.GetModelID(), request.GameName)
	log.Printf("Request context: Knowledge length=%d, Question length=%d", len(request.Knowledge), len(request.Question))

	bedrockRequest := &types.BedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []types.BedrockMessage{
			{
				Role:    "user",
				Content: userContent,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
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

func (b *BedrockProvider) extractTextFromResponse(response *types.BedrockResponse) (string, error) {
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
