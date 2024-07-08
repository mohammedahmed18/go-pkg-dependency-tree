GO_SOURCES := $(shell find . -type f -name "*.go")
RUN_ARGS :=

.PHONY: build
build:
	go build -o ./go-deps-graph $(GO_SOURCES)

.PHONY: run
run:
	./go-deps-graph $(RUN_ARGS)

.PHONY: all
all: build run

