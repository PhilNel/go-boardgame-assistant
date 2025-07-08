package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
)

type BedrockCreator struct {
	bedrockClient aws.BedrockClient
}

func NewBedrockCreator(bedrockClient aws.BedrockClient) *BedrockCreator {
	return &BedrockCreator{
		bedrockClient: bedrockClient,
	}
}

func (b *BedrockCreator) CreateEmbedding(ctx context.Context, text string) ([]float64, error) {
	request := &TitanRequest{
		InputText:  text,
		Dimensions: 256,
		Normalize:  true,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	response, err := b.bedrockClient.InvokeEmbeddingModel(ctx, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke embedding model: %w", err)
	}

	var embeddingResponse TitanResponse
	if err := json.Unmarshal(response, &embeddingResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding response: %w", err)
	}

	return embeddingResponse.Embedding, nil
}
