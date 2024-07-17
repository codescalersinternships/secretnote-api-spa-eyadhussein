build:
	@go build -o bin/noteserver ./cmd/main.go

run: build
	@./bin/noteserver

test:
	@go test -v ./...
