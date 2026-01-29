.PHONY: build install test vet clean run

BINARY := decision-tree-cli
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/decision-tree-cli

install:
	go install ./cmd/dt

test:
	go test ./... -v -count=1

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY)
