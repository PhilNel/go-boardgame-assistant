package knowledge

import (
	"context"
	"fmt"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
)

type S3Provider struct {
	s3Client *aws.S3Client
}

func NewS3Provider(cfg *config.S3) (*S3Provider, error) {
	s3Client, err := aws.NewS3Client(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &S3Provider{
		s3Client: s3Client,
	}, nil
}

func (s *S3Provider) GetFiles(ctx context.Context, gameName string) ([]string, error) {
	folder := fmt.Sprintf("games/%s/", strings.ToLower(gameName))
	return s.s3Client.ListObjectsWithPrefix(ctx, folder)
}

func (s *S3Provider) GetFileContent(ctx context.Context, filePath string) ([]byte, error) {
	return s.s3Client.GetObject(ctx, filePath)
}
