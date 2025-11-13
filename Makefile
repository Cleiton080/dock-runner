BIN_OUTPUT=./build

build-config-watcher:
	go build -o $(BIN_OUTPUT)/config-watcher ./cmd/config-watcher/main.go

build: build-config-watcher
