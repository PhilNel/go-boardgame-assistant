package main

import (
	"context"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/answer"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/handler"
	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
	"github.com/PhilNel/go-boardgame-assistant/internal/prompt"
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

	knowledgeProvider, err := knowledge.NewVectorProvider(cfg.DynamoDB, cfg.Bedrock, cfg.RAG)
	if err != nil {
		log.Fatalf("Failed to initialize vector knowledge provider: %v", err)
	}
	log.Printf("Vector knowledge provider initialized successfully")

	templateProvider := prompt.NewStaticTemplate()
	log.Printf("Template provider initialized successfully")

	answerProvider, err := answer.NewBedrockProvider(cfg.Bedrock, templateProvider)
	if err != nil {
		log.Fatalf("Failed to initialize answer provider: %v", err)
	}
	log.Printf("Answer provider initialized successfully")

	questionHandler = handler.NewQuestionHandler(knowledgeProvider, answerProvider)
	log.Printf("Question handler initialized successfully")
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

	log.Printf("Handler completed successfully with status: %d", response.StatusCode)
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
