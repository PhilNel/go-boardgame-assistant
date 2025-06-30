package aws

import (
	"context"
	"fmt"
	"time"

	configPkg "github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamoDBClient(cfg *configPkg.DynamoDB) (*DynamoDBClient, error) {
	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(awsCfg)

	return &DynamoDBClient{
		client: client,
	}, nil
}

func (d *DynamoDBClient) PutItem(ctx context.Context, tableName string, item interface{}) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      itemMap,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (d *DynamoDBClient) GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, result interface{}) error {
	output, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	if output.Item == nil {
		return fmt.Errorf("item not found")
	}

	err = attributevalue.UnmarshalMap(output.Item, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return nil
}

func (d *DynamoDBClient) Query(ctx context.Context, tableName string, indexName *string, keyCondition string, expressionValues map[string]types.AttributeValue, results interface{}) error {
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionValues,
	}

	if indexName != nil {
		input.IndexName = indexName
	}

	output, err := d.client.Query(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, results)
	if err != nil {
		return fmt.Errorf("failed to unmarshal results: %w", err)
	}

	return nil
}

func (d *DynamoDBClient) BatchWriteItems(ctx context.Context, tableName string, items []interface{}) error {
	const batchSize = 25 // DynamoDB batch write limit

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		writeRequests := make([]types.WriteRequest, 0, len(batch))

		for _, item := range batch {
			itemMap, err := attributevalue.MarshalMap(item)
			if err != nil {
				continue // Skip invalid items
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: itemMap,
				},
			})
		}

		if len(writeRequests) == 0 {
			continue
		}

		_, err := d.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to batch write items: %w", err)
		}
	}

	return nil
}

func (d *DynamoDBClient) UpdateItem(ctx context.Context, tableName string, key map[string]types.AttributeValue, updateExpression string, expressionValues map[string]types.AttributeValue) error {
	_, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	})
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
