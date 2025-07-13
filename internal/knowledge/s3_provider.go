package knowledge

import (
	"context"
	"fmt"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
)

type S3Provider struct {
	s3Client aws.S3Client
}

func NewS3Provider(s3Client aws.S3Client) *S3Provider {
	return &S3Provider{
		s3Client: s3Client,
	}
}

func (s *S3Provider) GetFiles(ctx context.Context, gameName string) ([]string, error) {
	folder := fmt.Sprintf("games/%s/", strings.ToLower(gameName))
	return s.s3Client.ListObjectsWithPrefix(ctx, folder)
}

func (s *S3Provider) GetFileContent(ctx context.Context, filePath string) ([]byte, error) {
	return s.s3Client.GetObject(ctx, filePath)
}
