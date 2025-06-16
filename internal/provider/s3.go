package provider

import (
	"context"
	"fmt"
	"strings"

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
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.config.Bucket),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var folders []string
	for _, prefix := range result.CommonPrefixes {
		if prefix.Prefix != nil {
			// Remove trailing slash
			folder := strings.TrimSuffix(*prefix.Prefix, "/")
			folders = append(folders, folder)
		}
	}

	return folders, nil
}
