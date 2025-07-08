package knowledge

import (
	"context"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/PhilNel/go-boardgame-assistant/internal/config"
	"github.com/PhilNel/go-boardgame-assistant/internal/utils"
)

type SearchStrategy interface {
	Search(ctx context.Context, chunks []*Chunk, query string, queryEmbedding []float64) ([]*SearchResult, error)
}

// HybridSearchStrategy combines vector and keyword search results
type HybridSearchStrategy struct {
	ragConfig     *config.RAG
	vectorWeight  float64
	keywordWeight float64
}

func NewHybridSearchStrategy(ragConfig *config.RAG, vectorWeight, keywordWeight float64) *HybridSearchStrategy {
	return &HybridSearchStrategy{
		ragConfig:     ragConfig,
		vectorWeight:  vectorWeight,
		keywordWeight: keywordWeight,
	}
}

func (h *HybridSearchStrategy) Search(ctx context.Context, chunks []*Chunk, query string, queryEmbedding []float64) ([]*SearchResult, error) {
	vectorResults := h.performVectorSearch(chunks, queryEmbedding)
	keywordResults := h.performKeywordSearch(chunks, query)

	combinedResults := h.combineResults(vectorResults, keywordResults)

	var filteredResults []*SearchResult
	minScore := h.ragConfig.MinSimilarity * 0.7

	for _, result := range combinedResults {
		if result.Similarity >= minScore {
			filteredResults = append(filteredResults, result)
		}
	}

	log.Printf("Hybrid search: vector=%d, keyword=%d, combined=%d, filtered=%d",
		len(vectorResults), len(keywordResults), len(combinedResults), len(filteredResults))

	return filteredResults, nil
}

func (h *HybridSearchStrategy) performVectorSearch(chunks []*Chunk, queryEmbedding []float64) []*SearchResult {
	var results []*SearchResult
	minSimilarity := h.ragConfig.MinSimilarity

	for _, chunk := range chunks {
		similarity := utils.CosineSimilarity(queryEmbedding, chunk.Embedding)

		if similarity >= minSimilarity {
			results = append(results, &SearchResult{
				Chunk:      chunk,
				Similarity: similarity,
			})
		}
	}

	log.Printf("Vector search found %d chunks above %.2f similarity", len(results), minSimilarity)
	return results
}

func (h *HybridSearchStrategy) performKeywordSearch(chunks []*Chunk, query string) []*SearchResult {
	queryTerms := h.tokenize(strings.ToLower(query))
	if len(queryTerms) == 0 {
		return []*SearchResult{}
	}

	var results []*SearchResult
	for _, chunk := range chunks {
		score := h.calculateTFIDFScore(queryTerms, chunk, chunks)
		if score > 0 {
			results = append(results, &SearchResult{
				Chunk:      chunk,
				Similarity: score,
			})
		}
	}

	log.Printf("Keyword search found %d chunks with TF-IDF scores > 0", len(results))
	return results
}

func (h *HybridSearchStrategy) combineResults(vectorResults, keywordResults []*SearchResult) []*SearchResult {
	scoreMap := make(map[string]*SearchResult)
	vectorScores := h.normalizeScores(vectorResults)
	for i, result := range vectorResults {
		scoreMap[result.Chunk.ID] = &SearchResult{
			Chunk:      result.Chunk,
			Similarity: vectorScores[i] * h.vectorWeight,
		}
	}

	keywordScores := h.normalizeScores(keywordResults)
	for i, result := range keywordResults {
		if existing, exists := scoreMap[result.Chunk.ID]; exists {
			existing.Similarity += keywordScores[i] * h.keywordWeight
		} else {
			scoreMap[result.Chunk.ID] = &SearchResult{
				Chunk:      result.Chunk,
				Similarity: keywordScores[i] * h.keywordWeight,
			}
		}
	}

	var results []*SearchResult
	for _, result := range scoreMap {
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	return results
}

func (h *HybridSearchStrategy) normalizeScores(results []*SearchResult) []float64 {
	if len(results) == 0 {
		return []float64{}
	}

	scores := make([]float64, len(results))
	var maxScore, minScore float64

	maxScore = results[0].Similarity
	minScore = results[0].Similarity

	for i, result := range results {
		scores[i] = result.Similarity
		if result.Similarity > maxScore {
			maxScore = result.Similarity
		}
		if result.Similarity < minScore {
			minScore = result.Similarity
		}
	}

	scoreRange := maxScore - minScore
	if scoreRange == 0 {
		// All scores are the same, return 1.0 for all
		for i := range scores {
			scores[i] = 1.0
		}
	} else {
		for i := range scores {
			scores[i] = (scores[i] - minScore) / scoreRange
		}
	}

	return scores
}

// Helper methods for keyword search
func (h *HybridSearchStrategy) tokenize(text string) []string {
	// Simple tokenization - split on whitespace and punctuation
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
	})

	// Filter out short words and common stop words
	var tokens []string
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true, "should": true,
	}

	for _, word := range words {
		word = strings.ToLower(word)
		if len(word) > 2 && !stopWords[word] {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

func (h *HybridSearchStrategy) calculateTFIDFScore(queryTerms []string, chunk *Chunk, allChunks []*Chunk) float64 {
	chunkTokens := h.tokenize(strings.ToLower(chunk.Content))
	if len(chunkTokens) == 0 {
		return 0
	}

	// Calculate term frequency for the chunk - How often a word appears in a document
	termFreq := make(map[string]int)
	for _, token := range chunkTokens {
		termFreq[token]++
	}

	// Calculate TF-IDF score
	// Inverse Document Frequency - Gives higher weight to rare/unique words
	var score float64
	for _, term := range queryTerms {
		tf := float64(termFreq[term]) / float64(len(chunkTokens))
		if tf > 0 {
			idf := h.calculateIDF(term, allChunks)
			score += tf * idf
		}
	}

	return score
}

func (h *HybridSearchStrategy) calculateIDF(term string, allChunks []*Chunk) float64 {
	docsWithTerm := 0
	for _, chunk := range allChunks {
		if strings.Contains(strings.ToLower(chunk.Content), term) {
			docsWithTerm++
		}
	}

	if docsWithTerm == 0 {
		return 0
	}

	return math.Log(float64(len(allChunks)) / float64(docsWithTerm))
}
