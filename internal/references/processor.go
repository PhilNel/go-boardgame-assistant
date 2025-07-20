package references

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
)

// Unicode superscript digits for footnotes
var footnoteNumbers = []string{"¹", "²", "³", "⁴", "⁵", "⁶", "⁷", "⁸", "⁹", "¹⁰", "¹¹", "¹²", "¹³", "¹⁴", "¹⁵", "¹⁶", "¹⁷", "¹⁸", "¹⁹", "²⁰"}

type ReferenceProcessor struct {
	referenceRepo ReferenceRepository
}

func NewReferenceProcessor(referenceRepo ReferenceRepository) *ReferenceProcessor {
	return &ReferenceProcessor{
		referenceRepo: referenceRepo,
	}
}

func (p *ReferenceProcessor) ProcessCitations(ctx context.Context, gameID, responseText string) (*ProcessedResponse, error) {
	log.Printf("Processing citations for game: %s, text length: %d", gameID, len(responseText))

	citations, err := p.extractCitations(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to extract citations: %w", err)
	}

	if len(citations) == 0 {
		log.Printf("No citations found in response text")
		return &ProcessedResponse{
			Response:   responseText,
			References: nil,
		}, nil
	}

	log.Printf("Found %d citations to process", len(citations))

	footnoteMap, references, err := p.buildFootnoteMapping(ctx, gameID, citations)
	if err != nil {
		return nil, fmt.Errorf("failed to build footnote mapping: %w", err)
	}

	processedText := p.replaceCitationsWithFootnotes(responseText, citations, footnoteMap)

	result := &ProcessedResponse{
		Response:   processedText,
		References: references,
	}

	log.Printf("Citation processing completed. Found %d unique references", len(references))
	return result, nil
}

func (p *ReferenceProcessor) extractCitations(text string) ([]*Citation, error) {
	// Regex pattern to match [[REFERENCE-ID]] or [[REFERENCE-ID,page]]
	pattern := `\[\[([A-Z0-9\-_]+)(?:,(\d+))?\]\]`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile citation regex: %w", err)
	}

	matches := regex.FindAllStringSubmatch(text, -1)
	matchIndices := regex.FindAllStringIndex(text, -1)

	var citations []*Citation
	for i, match := range matches {
		if len(match) < 2 {
			continue
		}

		citation := &Citation{
			Original:    match[0],
			ReferenceID: match[1],
			StartPos:    matchIndices[i][0],
			EndPos:      matchIndices[i][1],
		}

		var wasPageNumberCaptured = len(match) > 2 && match[2] != ""
		if wasPageNumberCaptured {
			citation.Page = match[2]
		}

		citations = append(citations, citation)
	}

	// Sort citations by position in text (for proper replacement)
	sort.Slice(citations, func(i, j int) bool {
		return citations[i].StartPos < citations[j].StartPos
	})

	return citations, nil
}

func (p *ReferenceProcessor) buildFootnoteMapping(ctx context.Context, gameID string, citations []*Citation) (map[string]string, []*ReferenceInfo, error) {
	footnoteMap := make(map[string]string)
	var references []*ReferenceInfo
	footnoteCounter := 1

	// Track unique reference IDs to avoid duplicate footnotes
	seenReferences := make(map[string]bool)

	for _, citation := range citations {
		citationKey := citation.ReferenceID
		if citation.Page != "" {
			citationKey += "," + citation.Page
		}

		// Skip if we've already processed this exact citation
		if seenReferences[citationKey] {
			continue
		}

		if footnoteCounter > len(footnoteNumbers) {
			log.Printf("WARNING: Too many unique citations (%d), using number fallback", footnoteCounter)
			footnoteMap[citationKey] = fmt.Sprintf("^%d", footnoteCounter)
		} else {
			footnoteMap[citationKey] = footnoteNumbers[footnoteCounter-1]
		}

		reference, err := p.referenceRepo.GetReference(ctx, gameID, citation.ReferenceID)
		if err != nil {
			log.Printf("WARNING: Failed to lookup reference %s: %v", citation.ReferenceID, err)
			// Create a placeholder reference for missing ones
			references = append(references, &ReferenceInfo{
				ID:      footnoteCounter,
				Title:   "[Reference not found]",
				Section: "",
				Page:    citation.Page,
				URL:     "",
			})
		} else {
			// Build page reference
			pageRef := reference.PageReference
			if citation.Page != "" {
				pageRef = fmt.Sprintf("p.%s", citation.Page)
			}

			references = append(references, &ReferenceInfo{
				ID:      footnoteCounter,
				Title:   reference.Title,
				Section: reference.Section,
				Page:    pageRef,
				URL:     reference.URL,
			})
		}

		seenReferences[citationKey] = true
		footnoteCounter++
	}

	return footnoteMap, references, nil
}

func (p *ReferenceProcessor) replaceCitationsWithFootnotes(text string, citations []*Citation, footnoteMap map[string]string) string {
	// Work backwards through citations to maintain correct positions
	result := text
	for i := len(citations) - 1; i >= 0; i-- {
		citation := citations[i]

		citationKey := citation.ReferenceID
		if citation.Page != "" {
			citationKey += "," + citation.Page
		}

		footnote, exists := footnoteMap[citationKey]
		if !exists {
			log.Printf("WARNING: No footnote found for citation %s", citationKey)
			continue
		}

		result = result[:citation.StartPos] + footnote + result[citation.EndPos:]
	}

	return result
}
