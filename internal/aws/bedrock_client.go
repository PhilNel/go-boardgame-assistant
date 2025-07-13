package aws

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

type AWSBedrockClient struct {
	client           *bedrockruntime.Client
	modelID          string
	embeddingModelID string
}

func NewAWSBedrockClient(config *config.Bedrock) (*AWSBedrockClient, error) {
	ctx := context.Background()

	log.Printf("Initializing Bedrock client with region: %s, model: %s, embedding model: %s",
		config.Region, config.ModelID, config.EmbeddingModelID)

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(awsCfg)

	return &AWSBedrockClient{
		client:           client,
		modelID:          config.ModelID,
		embeddingModelID: config.EmbeddingModelID,
	}, nil
}

func (b *AWSBedrockClient) InvokeModel(ctx context.Context, request *BedrockRequest) (*BedrockResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.modelID),
		ContentType: aws.String("application/json"),
		Body:        requestBody,
	}

	result, err := b.client.InvokeModel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	var response BedrockResponse
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (b *AWSBedrockClient) InvokeEmbeddingModel(ctx context.Context, requestBody []byte) ([]byte, error) {
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.embeddingModelID),
		ContentType: aws.String("application/json"),
		Body:        requestBody,
	}

	result, err := b.client.InvokeModel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke embedding model: %w", err)
	}

	return result.Body, nil
}

func (b *AWSBedrockClient) GetModelID() string {
	return b.modelID
}

func (b *AWSBedrockClient) GetEmbeddingModelID() string {
	return b.embeddingModelID
}
