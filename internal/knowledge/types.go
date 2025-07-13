package knowledge

import "context"

type Chunk struct {
	ID         string    `json:"id" dynamodbav:"chunk_id"`
	GameName   string    `json:"game_name" dynamodbav:"game_name"`
	SourceFile string    `json:"source_file" dynamodbav:"source_file"`
	Content    string    `json:"content" dynamodbav:"content"`
	Embedding  []float64 `json:"embedding" dynamodbav:"embedding"`
	TokenCount int       `json:"token_count" dynamodbav:"token_count"`
	CreatedAt  int64     `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt  int64     `json:"updated_at" dynamodbav:"updated_at"`
}

type SearchRequest struct {
	GameName      string  `json:"game_name"`
	Query         string  `json:"query"`
	MinSimilarity float64 `json:"min_similarity"`
	MaxTokens     int     `json:"max_tokens"`
	TopK          int     `json:"top_k"`
}

type SearchResult struct {
	Chunk      *Chunk  `json:"chunk"`
	Similarity float64 `json:"similarity"`
}

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

type ProcessingResult struct {
	JobID     string `json:"job_id"`
	GameName  string `json:"game_name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
}
