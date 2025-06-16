package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Provider struct {
	client *s3.Client
	config *config.S3
}

func NewS3Provider(config *config.S3) (*S3Provider, error) {
	ctx := context.Background()

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Provider{
		client: client,
		config: config,
	}, nil
}

func (s *S3Provider) ListFolders(ctx context.Context) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.config.Bucket),
		Delimiter: aws.String("/"),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var folders []string
	for _, prefix := range result.CommonPrefixes {
		if prefix.Prefix != nil {
			// Remove trailing slash
			folder := (*prefix.Prefix)[:len(*prefix.Prefix)-1]
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func (s *S3Provider) GetObject(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	// Read the entire object into memory
	return io.ReadAll(result.Body)
}

func (s *S3Provider) ListFilesInFolder(ctx context.Context, folder string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.Bucket),
		Prefix: aws.String(folder + "/"),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in folder: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		if obj.Key != nil && *obj.Key != folder+"/" {
			files = append(files, *obj.Key)
		}
	}

	return files, nil
}
