package knowledge

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
)

// When no knowledge chunks meet the similarity threshold
type NoRelevantKnowledgeError struct {
	GameName      string
	Query         string
	MinSimilarity float64
	ChunksFound   int
}

func (e *NoRelevantKnowledgeError) Error() string {
	return fmt.Sprintf("no chunks found above similarity threshold %.2f for query", e.MinSimilarity)
}

type VectorProvider struct {
	knowledgeRepo     KnowledgeRepository
	embeddingProvider EmbeddingProvider
	ragConfig         *config.RAG
	searchStrategy    SearchStrategy
}

func NewVectorProvider(knowledgeRepo KnowledgeRepository, embeddingProvider EmbeddingProvider, ragConfig *config.RAG) *VectorProvider {
	searchStrategy := NewHybridSearchStrategy(ragConfig, ragConfig.VectorWeight, ragConfig.KeywordWeight)

	return &VectorProvider{
		knowledgeRepo:     knowledgeRepo,
		embeddingProvider: embeddingProvider,
		ragConfig:         ragConfig,
		searchStrategy:    searchStrategy,
	}
}

func (v *VectorProvider) GetKnowledge(ctx context.Context, gameName string, query string) (string, error) {
	queryEmbedding, err := v.embeddingProvider.CreateEmbedding(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to create query embedding: %w", err)
	}

	chunks, err := v.knowledgeRepo.GetKnowledgeChunksByGame(ctx, gameName)
	if err != nil {
		return "", fmt.Errorf("failed to get knowledge chunks: %w", err)
	}

	log.Printf("Retrieved %d chunks for game '%s'", len(chunks), gameName)

	results, err := v.searchStrategy.Search(ctx, chunks, query, queryEmbedding)
	if err != nil {
		return "", fmt.Errorf("search strategy failed: %w", err)
	}

	if len(results) == 0 {
		return "", &NoRelevantKnowledgeError{
			GameName:      gameName,
			Query:         query,
			MinSimilarity: v.ragConfig.MinSimilarity,
			ChunksFound:   len(chunks),
		}
	}

	selectedResults := v.selectChunksWithinTokenBudget(results)
	combinedKnowledge := v.buildCombinedKnowledge(selectedResults, query)

	log.Printf("Search for '%s': found %d chunks, selected %d chunks with %d total tokens",
		query, len(results), len(selectedResults), v.calculateTotalTokens(selectedResults))

	return combinedKnowledge, nil
}

func (v *VectorProvider) selectChunksWithinTokenBudget(results []*SearchResult) []*SearchResult {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
	var selected []*SearchResult
	totalTokens := 0
	maxTokens := v.ragConfig.MaxTokens

	for _, result := range results {
		// Check if adding this chunk would exceed the token budget
		if totalTokens+result.Chunk.TokenCount <= maxTokens {
			selected = append(selected, result)
			totalTokens += result.Chunk.TokenCount
		} else {
			// Stop here - adding this chunk would exceed budget
			break
		}
	}

	log.Printf("Selected %d chunks with total tokens: %d (budget: %d)",
		len(selected), totalTokens, maxTokens)

	return selected
}

func (v *VectorProvider) buildCombinedKnowledge(selectedResults []*SearchResult, query string) string {
	log.Printf("=== SELECTED CHUNKS FOR QUERY: '%s' ===", query)

	var combinedKnowledge strings.Builder
	for i, result := range selectedResults {
		log.Printf("Chunk %d: File=%s, Tokens=%d, Score=%.4f",
			i+1, result.Chunk.SourceFile, result.Chunk.TokenCount, result.Similarity)
		combinedKnowledge.WriteString(fmt.Sprintf("Source %d (Score: %.2f, File: %s):\n",
			i+1, result.Similarity, result.Chunk.SourceFile))
		combinedKnowledge.WriteString(result.Chunk.Content)
		combinedKnowledge.WriteString("\n\n")
	}
	log.Printf("=== END SELECTED CHUNKS ===")

	return combinedKnowledge.String()
}

// calculateTotalTokens calculates the total token count for selected results
func (v *VectorProvider) calculateTotalTokens(selectedResults []*SearchResult) int {
	totalTokens := 0
	for _, result := range selectedResults {
		totalTokens += result.Chunk.TokenCount
	}
	return totalTokens
}
