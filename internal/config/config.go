package config

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Log      *Log
	S3       *S3
	Bedrock  *Bedrock
	DynamoDB *DynamoDB
	RAG      *RAG
	System   *System
}

type System struct {
	KnowledgeProvider string `long:"knowledge_provider" env:"KNOWLEDGE_PROVIDER" description:"Knowledge provider to use (s3 or vector)" default:"s3"`
}

type Bedrock struct {
	ModelID           string  `long:"bedrock_model_id" env:"BEDROCK_MODEL_ID" description:"Bedrock model ID to use" default:"anthropic.claude-3-haiku-20240307-v1:0"`
	EmbeddingModelID  string  `long:"bedrock_embedding_model_id" env:"BEDROCK_EMBEDDING_MODEL_ID" description:"Bedrock embedding model ID" default:"amazon.titan-embed-text-v2:0"`
	Region            string  `long:"aws_region_bedrock" env:"AWS_REGION" description:"AWS region to use" default:"eu-west-1"`
	AnswerMaxTokens   int     `long:"bedrock_max_tokens" env:"BEDROCK_ANSWER_MAX_TOKENS" description:"Maximum tokens to include in the answer" default:"1500"`
	AnswerTemperature float64 `long:"bedrock_temperature" env:"BEDROCK_ANSWER_TEMPERATURE" description:"Temperature for the Bedrock model answers" default:"0.1"`
	AnswerTopP        float64 `long:"bedrock_top_p" env:"BEDROCK_ANSWER_TOP_P" description:"TopP for the Bedrock model answers" default:"0.9"`
}

type Log struct {
	Level string `long:"log_level" env:"LOG_LEVEL" description:"Log level (debug, info, warn, error)" default:"info"`
}

type S3 struct {
	Bucket string `long:"knowledge_bucket" env:"KNOWLEDGE_BUCKET_NAME" description:"S3 bucket containing game knowledge files"`
	Region string `long:"aws_region_s3" env:"AWS_REGION" description:"AWS region to use" default:"eu-west-1"`
}

type DynamoDB struct {
	KnowledgeTable  string `long:"knowledge_table" env:"KNOWLEDGE_TABLE_NAME" description:"DynamoDB table for knowledge chunks"`
	JobsTable       string `long:"jobs_table" env:"JOBS_TABLE_NAME" description:"DynamoDB table for processing jobs"`
	FeedbackTable   string `long:"feedback_table" env:"FEEDBACK_TABLE_NAME" description:"DynamoDB table for feedback submissions"`
	ReferencesTable string `long:"references_table" env:"REFERENCES_TABLE_NAME" description:"DynamoDB table for game references"`
	Region          string `long:"aws_region_dynamodb" env:"AWS_REGION" description:"AWS region to use" default:"eu-west-1"`
}

type RAG struct {
	MinSimilarity  float64 `long:"rag_min_similarity" env:"RAG_MIN_SIMILARITY" description:"Minimum similarity threshold for vector search" default:"0.65"`
	MaxTokens      int     `long:"rag_max_tokens" env:"RAG_MAX_TOKENS" description:"Maximum tokens to include in context" default:"2000"`
	TopK           int     `long:"rag_top_k" env:"RAG_TOP_K" description:"Maximum number of chunks to retrieve" default:"10"`
	CacheTTLHours  int     `long:"cache_ttl_hours" env:"CACHE_TTL_HOURS" description:"Cache TTL in hours" default:"24"`
	MaxChunkTokens int     `long:"max_chunk_tokens" env:"MAX_CHUNK_TOKENS" description:"Maximum tokens per chunk" default:"500"`
	VectorWeight   float64 `long:"rag_vector_weight" env:"RAG_VECTOR_WEIGHT" description:"Weight for vector search in hybrid mode" default:"0.7"`
	KeywordWeight  float64 `long:"rag_keyword_weight" env:"RAG_KEYWORD_WEIGHT" description:"Weight for keyword search in hybrid mode" default:"0.3"`
}

func Load() (*Config, error) {
	opts := &Config{}
	_, err := flags.Parse(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return opts, nil
}
