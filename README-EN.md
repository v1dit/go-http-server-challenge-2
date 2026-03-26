# 🧠 Challenge 2

## Overview

This document describes the technical requirements and implementation architecture for **Challenge 2**, a continuation of the previous exercise with the addition of a new **business logic** layer.

The goal is to build an **In-Memory Login Service** — but, more important than what is built, is **how it is built**.

---

## 🧰 Tech Stack

| Component | Version |
|---|---|
| Language | Go 1.22+ |
| Router | `github.com/go-chi/chi/v5` |
| Hashing | `golang.org/x/crypto/bcrypt` *(optional, but recommended)* |
| Port | `:8080` |

---

## 🏗️ A New Way of Thinking — Separation of Concerns

In Challenge 1, all logic lived directly inside the *handlers*. That was a suitable approach for the context — routing, parameter parsing, simple responses. There was no real business logic to separate.

This challenge is different. From here on, there are **business decisions** to be made: is a *username* valid? Is a *password* correct? Does a session exist? These decisions **do not belong in the handler**.

The requested architecture rests on a clear and deliberate separation between two layers with distinct and non-transferable responsibilities:

**Handlers are solely responsible for the HTTP layer.** Their work begins and ends with the protocol: receive the request, parse the body, invoke the service, and translate the result into an HTTP response. A handler does not know what "a user already exists" means — it only knows that, when the service communicates that fact, it must respond with a specific status code. It contains no business conditions. It validates no domain rules. It is deliberately dumb.

**The `AuthService` is the sole owner of business logic.** It decides whether a *username* meets the minimum requirements, whether credentials are valid, whether a session token should be issued. Crucially, the service is **completely agnostic to the transport layer** — it does not know `http.Request`, does not write to `http.ResponseWriter`, and does not even know an HTTP server exists. It receives primitive data, applies rules, and returns results or domain errors.

This division is not architectural bureaucracy. It is what makes each layer **independently testable**, and what ensures that, in the future, the business logic can be reused in other contexts without any changes — whether a CLI, an async worker, or a gRPC endpoint.

> The handler does not decide — it executes. The service does not communicate — it decides.

---

## 1. The `AuthService` — Business Logic

The `AuthService` is the application's **source of truth**. It keeps registered users and active sessions in memory.

### Method `Register(username, password string) (User, error)`

**Business rules:**

| Condition | Behaviour |
|---|---|
| `username` shorter than 4 characters | Returns a validation error |
| `username` already exists | Returns `ErrUserExists` |
| Successful registration | Returns the created user |

---

### Method `Login(username, password string) (string, error)`

**Business rules:**

| Condition | Behaviour |
|---|---|
| User does not exist | Returns `ErrUnauthorized` |
| Incorrect password | Returns `ErrUnauthorized` |
| Successful authentication | Returns a session token |

---

### Domain Errors

The service must expose its own semantic errors. These errors allow handlers to map each failure scenario to the correct HTTP status code — without the handler needing any knowledge of the rules that caused them.

---

## 2. Initial Setup — Bootstrap

On application startup, an administrator user must be created for development and testing purposes:

| Field | Value |
|---|---|
| Username | `admin` |
| Password | `password123` |

---

## 3. Endpoints

### `POST /api/user/register`

Registers a new user in the system.

**Request body:**
```json
{
  "username": "intern1",
  "password": "safe-password"
}
```

**Success response (`201 Created`):**
```json
{
  "id": "a1b2c3d4",
  "username": "intern1",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Error responses:**

| Scenario | Code |
|---|---|
| `username` shorter than 4 characters | `400 Bad Request` |
| `username` already registered | `409 Conflict` |
| Invalid or missing JSON body | `400 Bad Request` |

---

### `POST /api/user/login`

Authenticates an existing user and issues a session token.

**Request body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Success response (`200 OK`):**
```json
{
  "token": "f7a3c9e1b2d4..."
}
```

**Error responses:**

| Scenario | Code |
|---|---|
| Invalid credentials | `401 Unauthorized` |
| Invalid or missing JSON body | `400 Bad Request` |

---

### `GET /api/user/profile`

Returns the authenticated user's profile based on the active session token.

**Required header:**
```
X-Session-Token: f7a3c9e1b2d4...
```

**Success response (`200 OK`):**
```json
{
  "id": "a1b2c3d4",
  "username": "admin",
  "created_at": "2025-01-15T09:00:00Z"
}
```

**Error responses:**

| Scenario | Code |
|---|---|
| Missing `X-Session-Token` header | `401 Unauthorized` |
| Invalid or non-existent token | `401 Unauthorized` |

---

## 🚀 How to Run

```bash
# Install dependencies
go mod tidy

# Start the server
go run main.go

# The server will be available at:
# http://localhost:8080
```

---

## 🧪 cURL Test Examples

```bash
# Register a new user
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "intern1", "password": "safe-password"}'

# Try to register a duplicate username (should return 409)
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "intern1", "password": "another-password"}'

# Try to register a username that is too short (should return 400)
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "ab", "password": "safe-password"}'

# Authenticate with the admin user
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'

# Authenticate with invalid credentials (should return 401)
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "wrongpassword"}'

# Get the authenticated user's profile (replace TOKEN with the value received on login)
curl http://localhost:8080/api/user/profile \
  -H "X-Session-Token: TOKEN_HERE"

# Try to access /profile without a token (should return 401)
curl http://localhost:8080/api/user/profile
```

---

## 📋 Endpoint Summary

| Method | Route | Auth | Description |
|---|---|---|---|
| `POST` | `/api/user/register` | ✗ | Register a new user |
| `POST` | `/api/user/login` | ✗ | Authenticate and issue a token |
| `GET` | `/api/user/profile` | ✅ `X-Session-Token` | Authenticated user's profile |
