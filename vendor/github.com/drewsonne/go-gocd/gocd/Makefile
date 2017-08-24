SHELL:=/bin/bash
.PHONY: test format

format:
	golint .
	gofmt -w -s .

test:
	go tool vet .
