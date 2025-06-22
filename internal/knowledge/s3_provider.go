package knowledge

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/aws"
	"github.com/PhilNel/go-boardgame-assistant/internal/config"
)

type S3Provider struct {
	s3Client *aws.S3Client
}

func NewS3Provider(config *config.S3) (*S3Provider, error) {
	s3Client, err := aws.NewS3Client(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	return &S3Provider{
		s3Client: s3Client,
	}, nil
}

func (p *S3Provider) GetKnowledge(ctx context.Context, gameName string) (string, error) {
	folder := strings.ToLower(gameName) + "/"

	files, err := p.s3Client.ListObjectsWithPrefix(ctx, folder)
	if err != nil {
		log.Printf("Failed to list files in folder: %v", err)
		return "", fmt.Errorf("failed to list game rule files: %w", err)
	}

	log.Printf("Retrieved %d files from S3 folder '%s': %v", len(files), folder, files)

	var combinedRules strings.Builder
	for _, file := range files {
		if p.isSupportedFile(file) {
			log.Printf("Processing file: %s", file)
			content, err := p.s3Client.GetObject(ctx, file)
			if err != nil {
				log.Printf("Failed to get file %s: %v", file, err)
				continue
			}
			combinedRules.WriteString(string(content))
			combinedRules.WriteString("\n\n")
		} else {
			log.Printf("Skipping file (unsupported extension): %s", file)
		}
	}

	if combinedRules.Len() == 0 {
		return "", fmt.Errorf("no game rule files found for game: %s", gameName)
	}

	return combinedRules.String(), nil
}

func (p *S3Provider) isSupportedFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt") || strings.HasSuffix(filename, ".md")
}
