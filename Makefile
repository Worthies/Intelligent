.PHONY: all build clean

BINARY_NAME=intelligent

all: build

build:
	go build -o $(BINARY_NAME) .

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 .

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o $(BINARY_NAME)-windows-arm64.exe .

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*
