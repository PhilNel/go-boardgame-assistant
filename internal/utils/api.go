package utils

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func CreateErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       fmt.Sprintf(`{"error": "%s"}`, message),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func CreateSuccessResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	responseBody, err := json.Marshal(data)
	if err != nil {
		return CreateErrorResponse(500, "Failed to process response"), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
