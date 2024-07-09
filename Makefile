GO_SOURCES := $(shell find . -type f -name "*.go")
RUN_ARGS :=

.PHONY: build
build:
	go build -o ./bin/go-deps-graph $(GO_SOURCES)

.PHONY: run
run:
	./bin/go-deps-graph $(RUN_ARGS)

.PHONY: all
all: build run

