package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/provider"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string   `json:"message"`
	Folders []string `json:"folders"`
}

var (
	cfg        *config.Config
	s3Provider *provider.S3Provider
)

func init() {
	cfg = config.Load()

	var err error
	s3Provider, err = provider.NewS3Provider(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize S3 provider: %v", err)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request data for request %s", request.RequestContext.RequestID)

	// List folders
	folders, err := s3Provider.ListFolders(ctx)
	if err != nil {
		log.Printf("Failed to list folders: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to list S3 folders",
		}, err
	}

	response := Response{
		Message: "Successfully listed S3 folders",
		Folders: folders,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
