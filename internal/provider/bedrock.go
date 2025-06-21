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
	systemPrompt := `You are an expert on Nemesis board game rules. Answer questions using ONLY the provided knowledge base, which contains the complete and accurate rules for Nemesis.

CRITICAL INSTRUCTIONS:
- Use only information from the provided knowledge base
- Do not use board game knowledge from your training data
- If the knowledge base doesn't contain enough information, say 'I don't have enough information about that specific rule'
- Always provide accurate information with proper citations [X, p.XX]
- Do not invent or assume rules that aren't explicitly stated
- If you don't know the answer, say 'I don't have enough information about that specific rule'

The provided knowledge base is authoritative and complete for Nemesis rules.`

	log.Printf("Using Bedrock model ID: %s", b.config.ModelID)

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
		ModelId:     aws.String(b.config.ModelID),
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
