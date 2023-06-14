export LOG_CODES_YAML=./log_codes.yaml

.PHONY: run-app run-terraform test clean

.PHONY: default
default: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run-app: ## Run the application
	@echo "Running the application..."
	go run ./cmd/app

run-coverage:#run-app ## Create the application
	cat ./app.log | go run ./cmd/coverage | tee coverage_report.txt

run-terraform: ## Run the application
	@echo "Running the application..."
	go run ./cmd/terraform

test: ## Test the application
	@echo "Testing the application..."
	@go test -v ./...

