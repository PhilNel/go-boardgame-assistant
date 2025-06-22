package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func getCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, x-api-key",
	}
}

func CreateErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	log.Printf("ERROR: Returning %d error: %s", statusCode, message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       fmt.Sprintf(`{"error": "%s"}`, message),
		Headers:    getCORSHeaders(),
	}
}

func CreateSuccessResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	responseBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: Failed to marshal response: %v", err)
		return CreateErrorResponse(500, "Failed to process response"), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    getCORSHeaders(),
	}, nil
}
