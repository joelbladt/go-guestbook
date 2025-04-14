# Makefile für das Go Gästebuch Projekt
SRC_PKG=github.com/joelbladt/go-guestbook/src/guestbook
TEST_DIR=./tests/Unit
COVER_OUT=./test-results/cover.out
COVER_HTML=./test-results/coverage.html

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

# Codequalität prüfen
lint:
	go vet ./...
	@if ! command -v golangci-lint >/dev/null; then \
		echo "🔍 golangci-lint not found. Install it for full linting."; \
	else \
		golangci-lint run; \
	fi

# Remove generated coverage files
clean:
	rm -f $(COVER_OUT) $(COVER_HTML)
