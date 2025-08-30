# Board Game Assistant Lambda Functions

This repository contains the Lambda functions for the Board Game Assistant project:

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
- Performs hybrid search using vector similarity and TFIDF scoring to find relevant rule sections
- Uses AWS Bedrock and Claude to generate contextual answers
- Builds citation list for answers and injects references into responses
- Returns natural language responses based on the game's rules

### 3. Feedback Handler (`feedback-handler`)
This Lambda function captures user feedback to improve answer quality and system performance.
- Stores user ratings and issue reports for each response
- Tracks conversation context to understand user intent
- Enables improvement of prompt engineering and knowledge retrieval

## Prerequisites

- Go 1.24 or later
- AWS CLI configured with appropriate credentials
- Terraform for infrastructure management

## Development

1. Install dependencies:
   ```bash
   make vendor
   ```

2. Each Lambda has a separate `build` command. For example, the `knowledge-processor` Lambda can be built using:
   ```bash
   make build-processor
   ```

3. Each Lambda can then be uploaded to the artefacts S3 bucket using the Deploy to AWS with the relevant `deploy` command:
   ```bash
   make deploy-processor
   ```
   The uploaded zip file can then be used to deploy the Lambda using the Terraform repository.

## Related Repositories

- [`knowledge-boardgame-assistant`](https://github.com/PhilNel/knowledge-boardgame-assistant) - Collection of structured board game rules in markdown format that forms the knowledge base for this project.

- [`infra-boardgame-assistant`](https://github.com/PhilNel/infra-boardgame-assistant) - Terraform configuration for deploying the infrastructure and managing Lambda permissions, S3 buckets, etc.

- [`vue-boardgame-assistant`](https://github.com/PhilNel/vue-boardgame-assistant) - The frontend Vue website that is used to interact with the Board Game Assistant functionality.