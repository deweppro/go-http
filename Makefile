SHELL=/bin/bash

gen:
	go generate ./...

test:
	go test -race -v ./...
