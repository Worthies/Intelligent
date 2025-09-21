# Intelligent

A simple HTTP server written in Go for handling ticket requests, with options for waiting for dynamic responses or proxying to another server.

## Features

- Handles requests to `/rpc/obtainTicket.action`
- Default mode: Responds with empty body and current date header
- `-wait` mode: Prints the request in cURL format, waits for a response file to be populated, and returns its content
- `-target` mode: Forwards requests to a specified URL, saves the response to a file, and returns the file path
- Configurable port and timeout

## Installation

Clone the repository and build:

```bash
git clone https://github.com/worthies/intelligent.git
cd intelligent
go build -o intelligent .
```

Or use the Makefile for cross-platform builds:

```bash
make build-linux-amd64
make build-windows-amd64
# etc.
```

## Usage

### Default Mode

```bash
./intelligent -port 8080
```

### Wait Mode

```bash
./intelligent -wait -port 8080
```

When a request is received, it prints the cURL command and the expected response file. Populate the file with the desired response content.

### Target Mode

```bash
./intelligent -target http://example.com -port 8080
```

Forwards requests to the target URL, saves responses to files, and returns the file paths.

### Options

- `-port`: Port to listen on (default: 80)
- `-wait`: Enable wait mode
- `-target`: Target URL for proxying
- `-timeout`: Timeout for wait mode (default: 180s)

## License

MIT License - see [LICENSE](LICENSE) for details.
