package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	claudeModelID = "anthropic.claude-3-haiku-20240307-v1:0"
)

type BedrockProvider struct {
	client *bedrockruntime.Client
	config *config.Bedrock
}

type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	Messages         []ClaudeMessage `json:"messages"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
}

type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewBedrockProvider(config *config.Bedrock) (*BedrockProvider, error) {
	ctx := context.Background()

	log.Printf("Initializing Bedrock provider with region: %s", config.Region)

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(awsCfg)

	return &BedrockProvider{
		client: client,
		config: config,
	}, nil
}

func (b *BedrockProvider) GenerateResponse(ctx context.Context, context string, question string) (string, error) {
	systemPrompt := `You are an expert board game assistant. Your role is to help players understand game rules and mechanics.
When answering questions:
1. Be clear and concise
2. Cite specific rules and page numbers when possible
3. If you're unsure about something, say so
4. Format your response in a way that's easy to read
5. Include relevant examples when helpful`

	log.Printf("Using Bedrock model ID: %s", claudeModelID)

	requestBody, err := json.Marshal(ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("%s\n\nGame Context:\n%s\n\nQuestion: %s", systemPrompt, context, question),
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(claudeModelID),
		ContentType: aws.String("application/json"),
		Body:        requestBody,
	}

	result, err := b.client.InvokeModel(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	var response ClaudeResponse
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	// Combine all text content
	var answer string
	for _, content := range response.Content {
		if content.Type == "text" {
			answer += content.Text
		}
	}

	return answer, nil
}
