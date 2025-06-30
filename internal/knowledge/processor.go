package knowledge

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
)

type KnowledgeRepository interface {
	SaveKnowledgeChunk(ctx context.Context, chunk *Chunk) error
	GetKnowledgeChunksByGame(ctx context.Context, gameName string) ([]*Chunk, error)
	BatchSaveKnowledgeChunks(ctx context.Context, chunks []*Chunk) error
}

type FileProvider interface {
	GetFiles(ctx context.Context, gameName string) ([]string, error)
	GetFileContent(ctx context.Context, filePath string) ([]byte, error)
}

type EmbeddingProvider interface {
	CreateEmbedding(ctx context.Context, text string) ([]float64, error)
}

type StatusRepository interface {
	CreateProcessingJob(ctx context.Context, gameName string, totalFiles int) (string, error)
	UpdateJobProgress(ctx context.Context, jobID string, progress int) error
	CompleteJob(ctx context.Context, jobID string, gameName string, processed, total int) error
	FailJob(ctx context.Context, jobID string, gameName string, errorMsg string) error
}

type Processor struct {
	fileProvider      FileProvider
	embeddingProvider EmbeddingProvider
	knowledgeRepo     KnowledgeRepository
	statusRepo        StatusRepository
	config            *config.RAG
}

func NewProcessor(fileProvider FileProvider, embeddingProvider EmbeddingProvider, knowledgeRepo KnowledgeRepository, statusRepo StatusRepository, cfg *config.RAG) *Processor {
	return &Processor{
		fileProvider:      fileProvider,
		embeddingProvider: embeddingProvider,
		knowledgeRepo:     knowledgeRepo,
		statusRepo:        statusRepo,
		config:            cfg,
	}
}

func (p *Processor) ProcessGame(ctx context.Context, gameName string) (*ProcessingResult, error) {
	log.Printf("Starting knowledge processing for game: %s", gameName)

	files, err := p.fileProvider.GetFiles(ctx, gameName)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}

	supportedFiles := p.filterSupportedFiles(files)
	if len(supportedFiles) == 0 {
		return nil, fmt.Errorf("no supported files found for game: %s", gameName)
	}

	// Create processing job with total count
	jobID, err := p.statusRepo.CreateProcessingJob(ctx, gameName, len(supportedFiles))
	if err != nil {
		return nil, fmt.Errorf("failed to create processing job: %w", err)
	}

	log.Printf("Processing %d files for game: %s", len(supportedFiles), gameName)

	var chunks []*Chunk
	processed := 0

	for _, file := range supportedFiles {
		log.Printf("Processing file: %s", file)

		content, err := p.fileProvider.GetFileContent(ctx, file)
		if err != nil {
			log.Printf("Failed to get file %s: %v", file, err)
			continue
		}

		chunk, err := p.createKnowledgeChunk(ctx, gameName, file, string(content))
		if err != nil {
			log.Printf("Failed to create chunk for file %s: %v", file, err)
			continue
		}

		chunks = append(chunks, chunk)
		processed++

		// Update progress periodically
		if processed%5 == 0 || processed == len(supportedFiles) {
			if err := p.statusRepo.UpdateJobProgress(ctx, jobID, processed); err != nil {
				log.Printf("Failed to update job progress: %v", err)
			}
		}
	}

	// Batch store chunks
	if len(chunks) > 0 {
		if err := p.knowledgeRepo.BatchSaveKnowledgeChunks(ctx, chunks); err != nil {
			// Fail the job and return error
			if failErr := p.statusRepo.FailJob(ctx, jobID, gameName, fmt.Sprintf("Failed to store chunks: %v", err)); failErr != nil {
				log.Printf("Failed to update job failure: %v", failErr)
			}
			return &ProcessingResult{
				JobID:   jobID,
				Status:  "failed",
				Message: fmt.Sprintf("Failed to store chunks: %v", err),
			}, fmt.Errorf("failed to store chunks: %w", err)
		}
	}

	if err := p.statusRepo.CompleteJob(ctx, jobID, gameName, processed, len(supportedFiles)); err != nil {
		log.Printf("Failed to update job completion: %v", err)
	}

	log.Printf("Knowledge processing completed for game: %s, processed: %d/%d",
		gameName, processed, len(supportedFiles))

	return &ProcessingResult{
		JobID:     jobID,
		GameName:  gameName,
		Status:    "completed",
		Message:   "Knowledge processing completed successfully",
		Processed: processed,
		Total:     len(supportedFiles),
	}, nil
}

func (p *Processor) createKnowledgeChunk(ctx context.Context, gameName, filePath, content string) (*Chunk, error) {
	// Count tokens (simple estimation: ~4 chars per token)
	tokenCount := len(content) / 4
	if tokenCount > p.config.MaxChunkTokens {
		log.Printf("Warning: chunk for %s exceeds max tokens (%d > %d)",
			filePath, tokenCount, p.config.MaxChunkTokens)
	}

	embedding, err := p.embeddingProvider.CreateEmbedding(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	chunkID := p.generateChunkID(gameName, filePath)

	chunk := &Chunk{
		ID:         chunkID,
		GameName:   gameName,
		SourceFile: filePath,
		Content:    content,
		Embedding:  embedding,
		TokenCount: tokenCount,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	return chunk, nil
}

func (p *Processor) filterSupportedFiles(files []string) []string {
	var supported []string
	for _, file := range files {
		if strings.HasSuffix(file, ".md") || strings.HasSuffix(file, ".txt") {
			supported = append(supported, file)
		}
	}
	return supported
}

func (p *Processor) generateChunkID(gameName, filePath string) string {
	combined := fmt.Sprintf("%s:%s", gameName, filePath)
	hash := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("%x", hash)
}

type ProcessingResult struct {
	JobID     string `json:"job_id"`
	GameName  string `json:"game_name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
}
