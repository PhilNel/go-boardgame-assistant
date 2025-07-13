package status

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoDBRepository struct {
	dynamoDB  aws.DynamoDBClient
	jobsTable string
}

func NewDynamoDBRepository(dynamoClient aws.DynamoDBClient, jobsTable string) *DynamoDBRepository {
	log.Printf("Initializing status repository with dynamoDB table: %s", jobsTable)

	return &DynamoDBRepository{
		dynamoDB:  dynamoClient,
		jobsTable: jobsTable,
	}
}

func (r *DynamoDBRepository) CreateProcessingJob(ctx context.Context, gameName string, totalFiles int) (string, error) {
	jobID := uuid.New().String()
	job := &Job{
		ID:        jobID,
		GameName:  gameName,
		Status:    "processing",
		Progress:  0,
		Total:     totalFiles,
		StartedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := r.SaveProcessingJob(ctx, job); err != nil {
		return "", fmt.Errorf("failed to create processing job: %w", err)
	}

	return jobID, nil
}

func (r *DynamoDBRepository) UpdateJobProgress(ctx context.Context, jobID string, progress int) error {
	key := map[string]dynamoTypes.AttributeValue{
		"id": &dynamoTypes.AttributeValueMemberS{Value: jobID},
	}

	updateExpression := "SET progress = :progress, updated_at = :updated_at"
	expressionValues := map[string]dynamoTypes.AttributeValue{
		":progress":   &dynamoTypes.AttributeValueMemberN{Value: fmt.Sprintf("%d", progress)},
		":updated_at": &dynamoTypes.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
	}

	return r.dynamoDB.UpdateItem(ctx, r.jobsTable, key, updateExpression, expressionValues)
}

func (r *DynamoDBRepository) CompleteJob(ctx context.Context, jobID string, gameName string, processed, total int) error {
	job := &Job{
		ID:        jobID,
		GameName:  gameName,
		Status:    "completed",
		Progress:  processed,
		Total:     total,
		UpdatedAt: time.Now().Unix(),
	}

	return r.SaveProcessingJob(ctx, job)
}

func (r *DynamoDBRepository) FailJob(ctx context.Context, jobID string, gameName string, errorMsg string) error {
	job := &Job{
		ID:        jobID,
		GameName:  gameName,
		Status:    "failed",
		UpdatedAt: time.Now().Unix(),
	}

	return r.SaveProcessingJob(ctx, job)
}

func (r *DynamoDBRepository) SaveProcessingJob(ctx context.Context, job *Job) error {
	return r.dynamoDB.PutItem(ctx, r.jobsTable, job)
}
