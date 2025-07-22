# Project settings
RULES_ASSISTANT_LAMBDA_CMD_DIR=cmd/question-handler
PROCESSOR_CMD_DIR=cmd/knowledge-processor
FEEDBACK_CMD_DIR=cmd/feedback-handler

BUCKET_NAME := boardgame-assistant-artefacts-dev-eu-west-1

# Lambda settings
RULES_ASSISTANT_LAMBDA_NAME := go-boardgame-rules-assistant
PROCESSOR_RULES_ASSISTANT_LAMBDA_NAME := go-boardgame-knowledge-processor
FEEDBACK_LAMBDA_NAME := go-boardgame-feedback-handler

BINARY_NAME := bootstrap

.PHONY: run-rules-assistant
run:
	go run $(RULES_ASSISTANT_LAMBDA_CMD_DIR)/main.go

.PHONY: run-processor
run-processor:
	go run $(PROCESSOR_CMD_DIR)/main.go

.PHONY: run-feedback
run-feedback:
	go run $(FEEDBACK_CMD_DIR)/main.go

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) ./$(RULES_ASSISTANT_LAMBDA_CMD_DIR)

.PHONY: build-processor
build-processor:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) ./$(PROCESSOR_CMD_DIR)

.PHONY: build-feedback
build-feedback:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) ./$(FEEDBACK_CMD_DIR)

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: package-rules-assistant
package-rules-assistant: build
	zip -j $(RULES_ASSISTANT_LAMBDA_NAME).zip $(BINARY_NAME)

.PHONY: package-processor
package-processor: build-processor
	zip -j $(PROCESSOR_LAMBDA_NAME).zip $(BINARY_NAME)

.PHONY: package-feedback
package-feedback: build-feedback
	zip -j $(FEEDBACK_LAMBDA_NAME).zip $(BINARY_NAME)

.PHONY: upload-rules-assistant
upload-rules-assistant:
	aws s3 cp $(RULES_ASSISTANT_LAMBDA_NAME).zip s3://$(BUCKET_NAME)/$(RULES_ASSISTANT_LAMBDA_NAME).zip

.PHONY: upload-processor
upload-processor:
	aws s3 cp $(PROCESSOR_LAMBDA_NAME).zip s3://$(BUCKET_NAME)/$(PROCESSOR_LAMBDA_NAME).zip

.PHONY: upload-feedback
upload-feedback:
	aws s3 cp $(FEEDBACK_LAMBDA_NAME).zip s3://$(BUCKET_NAME)/$(FEEDBACK_LAMBDA_NAME).zip

.PHONY: deploy-rules-assistant
deploy-rules-assistant: package-rules-assistant upload-rules-assistant

.PHONY: deploy-processor
deploy-processor: package-processor upload-processor

.PHONY: deploy-feedback
deploy-feedback: package-feedback upload-feedback

.PHONY: clean
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME) $(RULES_ASSISTANT_LAMBDA_NAME).zip $(PROCESSOR_LAMBDA_NAME).zip $(FEEDBACK_LAMBDA_NAME).zip 