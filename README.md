# Board Game Assistant Lambda Functions

This repository contains two Lambda functions for the Board Game Assistant project:

## Lambda Functions

### 1. Knowledge Processor (`knowledge-processor`)
This Lambda function processes and indexes board game rules from markdown files stored in S3.
- Reads game rule files from S3 storage
- Generates embeddings using AWS Bedrock
- Stores the processed knowledge chunks with embeddings in DynamoDB
- Tracks processing status for each game

### 2. Question Handler (`question-handler`)
This Lambda function serves as the main API backend for answering questions about board game rules.
- Receives user questions about specific board games
- Performs vector similarity search to find relevant rule sections
- Uses AWS Bedrock and Claude to generate contextual answers
- Returns natural language responses based on the game's rules

## Prerequisites

- Go 1.21 or later
- AWS CLI configured with appropriate credentials
- Terraform for infrastructure management

## Development

1. Install dependencies:
   ```bash
   make vendor
   ```

2. Build the Lambda package:
   ```bash
   make package
   ```

3. Deploy to AWS:
   ```bash
   make deploy
   ```

## Integration with Other Repositories

This Lambda function is part of a larger project with multiple repositories:

- `infra-boardgame-assistant/` - Terraform infrastructure
- `go-boardgame-assistant/` - This Lambda function
- `knowledge-boardgame-assistant/` - Game rules in markdown format

The infrastructure repository creates the necessary AWS resources and outputs values that this Lambda function uses. 