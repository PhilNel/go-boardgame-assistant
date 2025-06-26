package knowledge

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
