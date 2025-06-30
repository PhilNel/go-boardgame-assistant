package knowledge

import (
	"context"
	"fmt"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBRepository struct {
	dynamoDB       *aws.DynamoDBClient
	knowledgeTable string
}

func NewDynamoDBRepository(cfg *config.DynamoDB) (*DynamoDBRepository, error) {
	log.Printf("Initializing knowledge repository with dynamoDB table: %s", cfg.KnowledgeTable)

	dynamoDB, err := aws.NewDynamoDBClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create DynamoDB client: %w", err)
	}

	return &DynamoDBRepository{
		dynamoDB:       dynamoDB,
		knowledgeTable: cfg.KnowledgeTable,
	}, nil
}

func (r *DynamoDBRepository) SaveKnowledgeChunk(ctx context.Context, chunk *Chunk) error {
	err := r.dynamoDB.PutItem(ctx, r.knowledgeTable, chunk)
	if err != nil {
		return fmt.Errorf("failed to save knowledge chunk: %w", err)
	}

	log.Printf("Successfully stored knowledge chunk: %s", chunk.ID)
	return nil
}

func (r *DynamoDBRepository) GetKnowledgeChunksByGame(ctx context.Context, gameName string) ([]*Chunk, error) {
	var chunks []*Chunk

	err := r.dynamoDB.Query(ctx, r.knowledgeTable, nil,
		"game_name = :game_name",
		map[string]dynamoTypes.AttributeValue{
			":game_name": &dynamoTypes.AttributeValueMemberS{Value: gameName},
		}, &chunks)

	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge chunks: %w", err)
	}

	return chunks, nil
}

func (r *DynamoDBRepository) BatchSaveKnowledgeChunks(ctx context.Context, chunks []*Chunk) error {
	// Convert to []interface{} for generic batch write
	items := make([]interface{}, len(chunks))
	for i, chunk := range chunks {
		items[i] = chunk
	}

	err := r.dynamoDB.BatchWriteItems(ctx, r.knowledgeTable, items)
	if err != nil {
		return fmt.Errorf("failed to batch save knowledge chunks: %w", err)
	}

	log.Printf("Successfully batch saved %d knowledge chunks", len(chunks))
	return nil
}
