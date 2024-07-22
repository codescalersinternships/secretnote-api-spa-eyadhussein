BIN_DIR = bin
BIN_NAME = noteserver
CMD_PATH = ./cmd/main.go

.PHONY: build run test clean generate-swagger

build:
	@go build -o $(BIN_DIR)/$(BIN_NAME) $(CMD_PATH)

run: build
	@./$(BIN_DIR)/$(BIN_NAME) $(PORT)

test:
	@go test -v ./...

clean:
	@rm -rf $(BIN_DIR)/$(BIN_NAME)

generate-swagger:
	export PATH=$$PATH:$$HOME/go/bin
	bash -c "source ~/.bashrc && swag init -g $(CMD_PATH) -o ./docs"
