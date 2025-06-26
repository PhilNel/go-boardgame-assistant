package main

import (
	"context"
	"log"

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

	fileProvider, err := knowledge.NewS3Provider(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize file provider: %v", err)
	}

	embeddingProvider, err := embedding.NewBedrockCreator(cfg.Bedrock)
	if err != nil {
		log.Fatalf("Failed to initialize embedding provider: %v", err)
	}

	knowledgeRepo, err := knowledge.NewDynamoDBRepository(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Failed to initialize knowledge repository: %v", err)
	}

	statusRepo, err := status.NewDynamoDBRepository(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Failed to initialize status repository: %v", err)
	}

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
