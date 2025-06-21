package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(config *config.S3) (*S3Client, error) {
	ctx := context.Background()

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		client: client,
		bucket: config.Bucket,
	}, nil
}

func (s *S3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3Client) ListObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects with prefix: %w", err)
	}

	var keys []string
	for _, obj := range result.Contents {
		if obj.Key != nil && *obj.Key != prefix {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}
