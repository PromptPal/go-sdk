# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the PromptPal Golang SDK, a client library for interacting with the PromptPal service - a platform for managing and executing prompts. The SDK provides both synchronous and streaming execution capabilities.

## Build and Development Commands

```bash
# Run all tests
go test -v ./...

# Check for code issues
go vet ./...

# Format code
gofmt -w .

# Run specific test package
go test -v ./example/...

# Check dependencies
go mod tidy
go mod download
```

## Architecture Overview

### Core Components

1. **promptpal/client.go** - Main client implementation
   - `PromptPalClient` interface defines core methods: `Execute()` and `ExecuteStream()`
   - Client handles authentication, HTTP requests, and temporary token management
   - Built on top of `go-resty/resty/v2` HTTP client

2. **promptpal/http.go** - HTTP-related types
   - Request/response structs for API communication
   - Error handling structures

3. **promptpal/types.go** - Core type definitions
   - Configuration structures for code generation
   - Constants like `TEMPORARY_TOKEN_HEADER`

4. **example/** - Example usage and generated types
   - Contains test files demonstrating SDK usage
   - Includes generated types from PromptPal CLI

### Key Design Patterns

- The SDK uses functional options pattern for client configuration
- Supports both regular token authentication and temporary token authentication through callback
- Stream handling uses Server-Sent Events (SSE) format for real-time responses
- Error responses are properly typed and wrapped with context

### API Integration Points

- Base endpoint: `/api/v1/public/prompts/run/{pid}` for execution
- Stream endpoint: `/api/v1/public/prompts/run/{pid}/stream` for streaming
- Authentication via API token in `Authorization: API <token>` header
- Optional temporary token support via `X-TEMPORARY-TOKEN` header

## Testing Notes

Tests require a running PromptPal server at `localhost:7788` with valid tokens. The example tests will fail if the server is not available.

## Dependencies

- Go 1.20+
- github.com/go-resty/resty/v2 - HTTP client library