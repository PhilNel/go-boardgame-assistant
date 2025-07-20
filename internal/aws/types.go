package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient interface {
	PutItem(ctx context.Context, tableName string, item interface{}) error
	GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, result interface{}) error
	Query(ctx context.Context, tableName string, indexName *string, keyCondition string, expressionAttributeValues map[string]types.AttributeValue, result interface{}) error
	BatchWriteItems(ctx context.Context, tableName string, items []interface{}) error
	UpdateItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, updateExpression string, expressionValues map[string]types.AttributeValue) error
}

type S3Client interface {
	ListObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error)
	GetObject(ctx context.Context, key string) ([]byte, error)
}

type BedrockClient interface {
	InvokeModel(ctx context.Context, request *BedrockRequest) (*BedrockResponse, error)
	InvokeEmbeddingModel(ctx context.Context, requestBody []byte) ([]byte, error)
	GetModelID() string
	GetEmbeddingModelID() string
}

type BedrockRequest struct {
	AnthropicVersion string           `json:"anthropic_version"`
	Messages         []BedrockMessage `json:"messages"`
	MaxTokens        int              `json:"max_tokens,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
}

type BedrockMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BedrockResponse struct {
	Content []BedrockContent `json:"content"`
}

type BedrockContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
