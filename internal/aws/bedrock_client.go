package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient struct {
	client  *bedrockruntime.Client
	modelID string
}

func NewBedrockClient(config *config.Bedrock) (*BedrockClient, error) {
	ctx := context.Background()

	log.Printf("Initializing Bedrock client with region: %s, model: %s", config.Region, config.ModelID)

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(awsCfg)

	return &BedrockClient{
		client:  client,
		modelID: config.ModelID,
	}, nil
}

func (b *BedrockClient) InvokeModel(ctx context.Context, request *types.BedrockRequest) (*types.BedrockResponse, error) {
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

	var response types.BedrockResponse
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (b *BedrockClient) GetModelID() string {
	return b.modelID
}
