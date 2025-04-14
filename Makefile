# Makefile für das Go Gästebuch Projekt
APP_NAME=guestbook
SRC_PKG=github.com/joelbladt/go-guestbook/src/guestbook
TEST_DIR=./tests/Unit
BIN_DIR=./bin
COVER_OUT=./test-results/cover.out
COVER_HTML=./test-results/coverage.html

# Default Task
all:
	test
	build

# Start development server with live reloading using Air
dev:
	@if ! command -v air >/dev/null; then \
		echo "🔍 Air not found. Install it to starting a dev server."; \
	else \
		echo "🚀 Starting dev server with Air..."; \
		air; \
	fi

# Run tests with coverage and open HTML report
test:
	mkdir -p ./test-results
	go test -coverprofile=$(COVER_OUT) -coverpkg=$(SRC_PKG) $(TEST_DIR)
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	open $(COVER_HTML)

# Run only tests (no coverage)
test-only:
	go test -v $(TEST_DIR)

# Show coverage in terminal
coverage:
	go test -cover -coverpkg=$(SRC_PKG) $(TEST_DIR)

# Check Codequality
lint:
	go vet ./...
	@if ! command -v golangci-lint >/dev/null; then \
		echo "🔍 golangci-lint not found. Install it for full linting."; \
	else \
		golangci-lint run; \
	fi

# Build binary
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) .

# Multiplattform-Release (Linux, macOS, Windows)
release:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64   go build -o $(BIN_DIR)/$(APP_NAME)-linux .
	GOOS=darwin GOARCH=amd64  go build -o $(BIN_DIR)/$(APP_NAME)-macos .
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME).exe .

# Remove generated coverage files
clean:
	rm -f $(COVER_OUT) $(COVER_HTML)
