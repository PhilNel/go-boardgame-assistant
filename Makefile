# Project settings
CMD_DIR=cmd/lambda
BIN_NAME=boardgame-assistant
OUTPUT_BIN=./bin/$(BIN_NAME)
BUCKET_NAME := boardgame-assistant-artefacts-dev-eu-west-1
LAMBDA_NAME := go-boardgame-rules-assistant
ZIP_FILE := $(LAMBDA_NAME).zip
BINARY_NAME := bootstrap

.PHONY: run
run:
	go run $(CMD_DIR)/main.go

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) ./$(CMD_DIR)

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

.PHONY: package
package: build
	zip -j $(ZIP_FILE) $(BINARY_NAME)

.PHONY: upload
upload:
	aws s3 cp $(ZIP_FILE) s3://$(BUCKET_NAME)/$(ZIP_FILE)

.PHONY: deploy
deploy: package upload

.PHONY: clean
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME) $(ZIP_FILE) 