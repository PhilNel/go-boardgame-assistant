package references

import (
	"context"
	"fmt"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBRepository struct {
	dynamoDB        aws.DynamoDBClient
	referencesTable string
}

func NewDynamoDBRepository(dynamoClient aws.DynamoDBClient, referencesTable string) *DynamoDBRepository {
	log.Printf("Initializing references repository with DynamoDB table: %s", referencesTable)

	return &DynamoDBRepository{
		dynamoDB:        dynamoClient,
		referencesTable: referencesTable,
	}
}

func (r *DynamoDBRepository) GetReference(ctx context.Context, gameID, referenceID string) (*Reference, error) {
	key := map[string]dynamoTypes.AttributeValue{
		"gameId":      &dynamoTypes.AttributeValueMemberS{Value: gameID},
		"referenceId": &dynamoTypes.AttributeValueMemberS{Value: referenceID},
	}

	var reference Reference
	err := r.dynamoDB.GetItem(ctx, r.referencesTable, key, &reference)
	if err != nil {
		return nil, fmt.Errorf("failed to get reference %s for game %s: %w", referenceID, gameID, err)
	}

	log.Printf("Successfully retrieved reference: %s-%s", gameID, referenceID)
	return &reference, nil
}
