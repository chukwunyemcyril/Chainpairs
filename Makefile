# Makefile for Chainpairs (Go crypto exchange project)

BINARY_NAME=main

## Build the project
build:
	go build -buildvcs=false -o $(BINARY_NAME) .

## Run the project (builds first if needed)
run: build
	./$(BINARY_NAME)

## Clean up build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
