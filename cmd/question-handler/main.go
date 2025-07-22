package main

import (
	"context"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/answer"
	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/embedding"
	"github.com/PhilNel/go-boardgame-assistant/internal/handler"
	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
	"github.com/PhilNel/go-boardgame-assistant/internal/prompt"
	"github.com/PhilNel/go-boardgame-assistant/internal/references"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var questionHandler *handler.QuestionHandler

func init() {
	log.Printf("Starting Lambda initialization")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v", cfg)

	dynamoClient, err := aws.NewDynamoDBClient(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}
	knowledgeRepo := knowledge.NewDynamoDBRepository(dynamoClient, cfg.DynamoDB.KnowledgeTable)

	bedrockClient, err := aws.NewAWSBedrockClient(cfg.Bedrock)
	if err != nil {
		log.Fatalf("Failed to create Bedrock client: %v", err)
	}
	embeddingProvider := embedding.NewBedrockCreator(bedrockClient)

	templateProvider := prompt.NewStaticTemplate()

	referencesRepo := references.NewDynamoDBRepository(dynamoClient, cfg.DynamoDB.ReferencesTable)
	referenceProcessor := references.NewReferenceProcessor(referencesRepo)

	answerProvider := answer.NewBedrockProvider(bedrockClient, templateProvider, cfg.Bedrock)
	knowledgeProvider := knowledge.NewVectorProvider(knowledgeRepo, embeddingProvider, cfg.RAG)
	questionHandler = handler.NewQuestionHandler(knowledgeProvider, answerProvider, referenceProcessor)

	log.Printf("Lambda initialized successfully with references support")
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: Lambda handler panicked: %v", r)
		}
	}()

	response, err := questionHandler.Handle(ctx, request)
	if err != nil {
		log.Printf("ERROR: Handler returned error: %v", err)
		return utils.CreateErrorResponse(500, "Internal server error"), nil
	}

	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
