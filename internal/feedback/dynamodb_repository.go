package feedback

import (
	"context"
	"fmt"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
)

type DynamoDBRepository struct {
	dynamoDB      aws.DynamoDBClient
	feedbackTable string
}

func NewDynamoDBRepository(dynamoClient aws.DynamoDBClient, feedbackTable string) *DynamoDBRepository {
	log.Printf("Initializing feedback repository with dynamoDB table: %s", feedbackTable)

	return &DynamoDBRepository{
		dynamoDB:      dynamoClient,
		feedbackTable: feedbackTable,
	}
}

func (r *DynamoDBRepository) SaveFeedback(ctx context.Context, feedback *FeedbackRecord) error {
	err := r.dynamoDB.PutItem(ctx, r.feedbackTable, feedback)
	if err != nil {
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	log.Printf("Successfully stored feedback: %s for message: %s", feedback.FeedbackID, feedback.MessageID)
	return nil
}
