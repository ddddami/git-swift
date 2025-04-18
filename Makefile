BINARY_NAME=gitch
BUILD_DIR=bin

.PHONY: build clean run rebuild
build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

clean:
	@rm -rf $(BUILD_DIR)
	@echo "✅ Cleaned."

run:
	@go run ./cmd

rebuild: clean build

