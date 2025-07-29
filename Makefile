APP_NAME=adowork
SRC_DIR=src
BIN_DIR=bin
APP_PATH=$(BIN_DIR)/$(APP_NAME)
MAIN=$(SRC_DIR)/main.go

.PHONY: all build test clean run

all: build

build:
	mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && go build -o ../$(APP_PATH) .

test:
	cd $(SRC_DIR) && go test ./...

clean:
	rm -f $(APP_PATH)

run:
	cd $(SRC_DIR) && go run .

