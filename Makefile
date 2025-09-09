build:
	go build -o chainpairs main.go

run: build 
	go run main.go &&./chainpairs 

clean:
	rm -f chainpairs

test:
	go test -v ./...

.PHONY: build run clean test