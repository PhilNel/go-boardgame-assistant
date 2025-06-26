package aws

type BedrockRequest struct {
	AnthropicVersion string           `json:"anthropic_version"`
	Messages         []BedrockMessage `json:"messages"`
	MaxTokens        int              `json:"max_tokens,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
}

type BedrockMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BedrockResponse struct {
	Content []BedrockContent `json:"content"`
}

type BedrockContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
