package main

import (
	"context"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/feedback"
	"github.com/PhilNel/go-boardgame-assistant/internal/handler"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var feedbackHandler *handler.FeedbackHandler

func init() {
	log.Printf("Starting Feedback Handler Lambda initialization")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dynamoClient, err := aws.NewDynamoDBClient(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	feedbackRepo := feedback.NewDynamoDBRepository(dynamoClient, cfg.DynamoDB.FeedbackTable)

	feedbackBusinessHandler := feedback.NewHandler(feedbackRepo)
	feedbackHandler = handler.NewFeedbackHandler(feedbackBusinessHandler)

	log.Printf("Feedback Handler Lambda initialized successfully")
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: Lambda handler panicked: %v", r)
		}
	}()

	response, err := feedbackHandler.Handle(ctx, request)
	if err != nil {
		log.Printf("ERROR: Handler returned error: %v", err)
		return utils.CreateErrorResponse(500, "Internal server error"), nil
	}

	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
