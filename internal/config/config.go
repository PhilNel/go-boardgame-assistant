package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Log     *Log
	Runtime *Runtime
	S3      *S3
}

type Runtime struct {
	LambdaRuntimeAPI string `long:"lambda_runtime_api" description:"Set automatically by AWS Lambda" env:"AWS_LAMBDA_RUNTIME_API" default:""`
}

type Log struct {
	Level string `long:"log_level" env:"LOG_LEVEL" description:"Log level (debug, info, warn, error)" default:"info"`
}

type S3 struct {
	Bucket string `long:"knowledge_bucket" env:"KNOWLEDGE_BUCKET_NAME" description:"S3 bucket containing game knowledge files" required:"true"`
	Region string `long:"aws_region" env:"AWS_REGION" description:"AWS region to use" default:"us-east-1"`
}

var parsed *Config

func Load() *Config {
	if parsed != nil {
		return parsed
	}

	opts := &Config{}
	_, err := flags.Parse(opts)
	if err != nil {
		os.Exit(1) // flags already prints errors/help
	}

	parsed = opts
	return opts
}
