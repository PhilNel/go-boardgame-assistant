package config

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Log     *Log
	S3      *S3
	Bedrock *Bedrock
}

type Bedrock struct {
	Region string `long:"aws_region_bedrock" env:"AWS_REGION" description:"AWS region to use" default:"us-east-1"`
}

type Log struct {
	Level string `long:"log_level" env:"LOG_LEVEL" description:"Log level (debug, info, warn, error)" default:"info"`
}

type S3 struct {
	Bucket string `long:"knowledge_bucket" env:"KNOWLEDGE_BUCKET_NAME" description:"S3 bucket containing game knowledge files" required:"true"`
	Region string `long:"aws_region_s3" env:"AWS_REGION" description:"AWS region to use" default:"us-east-1"`
}

func Load() (*Config, error) {
	opts := &Config{}
	_, err := flags.Parse(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return opts, nil
}
