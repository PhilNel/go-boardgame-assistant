package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/provider"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	GameName string `json:"gameName"`
	Question string `json:"question"`
}

type Response struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

var (
	s3Provider      *provider.S3Provider
	bedrockProvider *provider.BedrockProvider
)

func init() {
	log.Printf("Starting Lambda initialization")

	// Initialize S3 provider
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v", cfg)

	s3Provider, err = provider.NewS3Provider(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize S3 provider: %v", err)
	}
	log.Printf("S3 provider initialized successfully")

	bedrockProvider, err = provider.NewBedrockProvider(cfg.Bedrock)
	if err != nil {
		log.Fatalf("Failed to initialize Bedrock provider: %v", err)
	}
	log.Printf("Bedrock provider initialized successfully")
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid request format: %v"}`, err),
		}, nil
	}

	if req.GameName == "" || req.Question == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "gameName and question are required"}`,
		}, nil
	}

	// Get all game rules files from S3
	folder := strings.ToLower(req.GameName)
	files, err := s3Provider.ListFilesInFolder(ctx, folder)
	if err != nil {
		log.Printf("Failed to list files in folder: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to list game rule files"}`,
		}, nil
	}

	var combinedRules strings.Builder
	for _, file := range files {
		if strings.HasSuffix(file, ".txt") {
			content, err := s3Provider.GetObject(ctx, file)
			if err != nil {
				log.Printf("Failed to get file %s: %v", file, err)
				continue
			}
			combinedRules.WriteString(string(content))
			combinedRules.WriteString("\n\n")
		}
	}

	if combinedRules.Len() == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "No game rule files found"}`,
		}, nil
	}

	// Generate response using Bedrock
	answer, err := bedrockProvider.GenerateResponse(ctx, combinedRules.String(), req.Question)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to generate response"}`,
		}, nil
	}

	response := Response{
		Answer: answer,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to process response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
