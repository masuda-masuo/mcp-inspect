.PHONY: build build-mocks test demo clean

BINARY    := mcp-inspect
CMD       := ./cmd/inspect

build:
	go build -o $(BINARY) $(CMD)

build-mocks:
	go build -o testdata/mock-servers/safe/mock-server       ./testdata/mock-servers/safe
	go build -o testdata/mock-servers/destructive/mock-server  ./testdata/mock-servers/destructive
	go build -o testdata/mock-servers/hidden-tools/mock-server ./testdata/mock-servers/hidden-tools

# Run all Go tests
test:
	go test ./...

# Full demo: build everything → run inspect on the demo config
demo: build build-mocks
	./$(BINARY) --config testdata/demo-config.json --output demo-report.html

# JSON demo (CI-friendly)
demo-json: build build-mocks
	./$(BINARY) --config testdata/demo-config.json --format json

clean:
	rm -f $(BINARY) demo-report.html mcp-report.html
	rm -f testdata/mock-servers/safe/mock-server
	rm -f testdata/mock-servers/destructive/mock-server
	rm -f testdata/mock-servers/hidden-tools/mock-server
