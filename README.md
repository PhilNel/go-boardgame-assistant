# Board Game Assistant Lambda

This repository contains the Lambda function for the Board Game Assistant project. The Lambda function serves as the backend for answering questions about board game rules using AWS Bedrock and Claude.

## Prerequisites

- Go 1.21 or later
- AWS CLI configured with appropriate credentials
- Terraform for infrastructure management

## Development

1. Install dependencies:
   ```bash
   make vendor
   ```

2. Run locally:
   ```bash
   make run
   ```

3. Build the Lambda package:
   ```bash
   make package
   ```

4. Deploy to AWS:
   ```bash
   make deploy
   ```

## Environment Variables

The Lambda function expects the following environment variables:

- `KNOWLEDGE_BUCKET_NAME` - S3 bucket containing game knowledge files (required)
- `AWS_REGION` - AWS region for both S3 and Bedrock services (default: us-east-1)
- `BEDROCK_MODEL_ID` - Bedrock model ID to use (default: anthropic.claude-3-haiku-20240307-v1:0)
- `LOG_LEVEL` - Logging level (debug, info, warn, error) (default: info)

## Integration with Other Repositories

This Lambda function is part of a larger project with multiple repositories:

- `infra-boardgame-assistant/` - Terraform infrastructure
- `go-boardgame-assistant/` - This Lambda function
- `knowledge-boardgame-assistant/` - Game rules in markdown format

The infrastructure repository creates the necessary AWS resources and outputs values that this Lambda function uses. 