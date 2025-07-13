package main

import (
	"context"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/embedding"
	"github.com/PhilNel/go-boardgame-assistant/internal/handler"
	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
	"github.com/PhilNel/go-boardgame-assistant/internal/status"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var processingHandler *handler.ProcessingHandler

func init() {
	log.Printf("Starting Knowledge Processor Lambda initialization")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	s3Client, err := aws.NewS3Client(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	fileProvider := knowledge.NewS3Provider(s3Client)

	bedrockClient, err := aws.NewAWSBedrockClient(cfg.Bedrock)
	if err != nil {
		log.Fatalf("Failed to create Bedrock client: %v", err)
	}
	embeddingProvider := embedding.NewBedrockCreator(bedrockClient)

	dynamoClient, err := aws.NewDynamoDBClient(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	knowledgeRepo := knowledge.NewDynamoDBRepository(dynamoClient, cfg.DynamoDB.KnowledgeTable)
	statusRepo := status.NewDynamoDBRepository(dynamoClient, cfg.DynamoDB.JobsTable)

	processor := knowledge.NewProcessor(fileProvider, embeddingProvider, knowledgeRepo, statusRepo, cfg.RAG)
	processingHandler = handler.NewProcessingHandler(processor)

	log.Printf("Knowledge Processor Lambda initialized successfully")
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: Lambda handler panicked: %v", r)
		}
	}()

	response, err := processingHandler.Handle(ctx, request)
	if err != nil {
		log.Printf("ERROR: Handler returned error: %v", err)
		return utils.CreateErrorResponse(500, "Internal server error"), nil
	}

	log.Printf("Handler completed successfully with status: %d", response.StatusCode)
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
