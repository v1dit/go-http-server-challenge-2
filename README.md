# Go HTTP Server (challenge 2)

## Overview
this repo contains a simple HTTP server built in Go using the chi router, implementing a basic authentication system with in-memory storage

## Features

### Routes
All endpoints are under the `/api/user` prefix:

- `POST /api/user/register` → register a new user
- `POST /api/user/login` → authenticate a user and return a session token
- `GET /api/user/profile` → get the authenticated user's profile (requires session token)

### Authentication
- in-memory user storage
- in-memory session/token management
- default user:
  - username: `admin`
  - password: `password123`
- profile endpoint requires `X-Session-Token` header

### Architecture
- `main.go` → initializes chi router and starts server on `:8080`
- `internal/handlers` → handles HTTP requests/responses
- `internal/service` → contains core authentication logic (register, login, profile)

this separation keeps HTTP logic and business logic clean and testable

### Changes
Refactor is complete: session validation is now middleware-based, and ProfileHandler is focused on business logic only.
* Added AuthMiddleware in internal/middleware/auth_middleware.go:
* Reads X-Session-Token
* Validates token 
* Returns consistent 401/500 responses on auth failure
* Stores authenticated user in request context
* Updated ProfileHandler 
* Updated route wiring in main.go:

### Verification
* Ran gofmt on updated files
* Ran go test ./... successfully
* Checked lints on edited files: no issues
Follows same arhcietructal rirection from challenge 1

## Run Locally

```bash
go mod tidy
go run main.go
