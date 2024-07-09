GO_SOURCES := $(shell find . -path ./test_samples -prune -o -type f -name "*.go" -print)
RUN_ARGS :=

.PHONY: build
build:
	go build -o ./bin/go-deps-graph $(GO_SOURCES)

.PHONY: run
run:
	./bin/go-deps-graph $(RUN_ARGS)

.PHONY: all
all: build run

