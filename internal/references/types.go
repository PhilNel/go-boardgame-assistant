package references

import "context"

type Reference struct {
	GameID        string `json:"gameId" dynamodbav:"gameId"`
	ReferenceID   string `json:"referenceId" dynamodbav:"referenceId"`
	Type          string `json:"type" dynamodbav:"type"`
	Title         string `json:"title" dynamodbav:"title"`
	Section       string `json:"section" dynamodbav:"section"`
	PageReference string `json:"pageReference" dynamodbav:"pageReference"`
	URL           string `json:"url" dynamodbav:"url"`
}

type ProcessedResponse struct {
	Response   string           `json:"response"`
	References []*ReferenceInfo `json:"references,omitempty"`
}

type ReferenceInfo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Section string `json:"section"`
	Page    string `json:"page"`
	URL     string `json:"url"`
}

type Citation struct {
	Original    string // Original citation text like [[R1-SLIME,17]]
	ReferenceID string // R1-SLIME
	Page        string // 17 (optional)
	StartPos    int    // Position in text
	EndPos      int    // End position in text
}

type ReferenceRepository interface {
	GetReference(ctx context.Context, gameID, referenceID string) (*Reference, error)
}

type Processor interface {
	Process(ctx context.Context, gameID, responseText string) (*ProcessedResponse, error)
}
