package main

import (
	"context"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/answer"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/handler"
	"github.com/PhilNel/go-boardgame-assistant/internal/knowledge"
	"github.com/PhilNel/go-boardgame-assistant/internal/prompt"
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

	knowledgeProvider, err := knowledge.NewS3Provider(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize knowledge provider: %v", err)
	}
	log.Printf("Knowledge provider initialized successfully")

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
	return questionHandler.Handle(ctx, request)
}

func main() {
	lambda.Start(handleRequest)
}
