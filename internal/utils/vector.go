package utils

import (
	"log"
	"math"
)

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where 1 indicates identical vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		log.Printf("DEBUG: Similarity calculation failed - length mismatch: a=%d, b=%d", len(a), len(b))
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		log.Printf("DEBUG: Similarity calculation failed - zero norm: normA=%.6f, normB=%.6f", normA, normB)
		return 0.0
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))

	return similarity
}
