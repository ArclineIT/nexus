# Nexus Control Panel

**Your central command. One login, every tool.**

Nexus is the identity hub for the Arcline platform. It acts as the single source of truth for users, roles, sessions, and SSO integrations — connecting Portal, Billing, Git, Monitoring, and future tools through one unified login.

## Quick Start

```bash
make run
```

The server starts on `http://0.0.0.0:8080` by default.

## Endpoints

| Method | Path            | Auth     | Description              |
|--------|-----------------|----------|--------------------------|
| GET    | `/health`       | Public   | Health check             |
| GET    | `/ready`        | Public   | Readiness probe          |
| POST   | `/auth/login`   | Public   | Authenticate, get tokens |
| POST   | `/auth/refresh` | Public   | Refresh access token     |
| GET    | `/auth/me`      | Bearer   | Current user info        |

## Configuration

All config is via environment variables:

| Variable                | Default                 | Description               |
|-------------------------|-------------------------|---------------------------|
| `NEXUS_HOST`            | `0.0.0.0`               | Server host               |
| `NEXUS_PORT`            | `8080`                  | Server port               |
| `NEXUS_DB_HOST`         | `localhost`             | PostgreSQL host           |
| `NEXUS_DB_PORT`         | `5432`                  | PostgreSQL port           |
| `NEXUS_DB_USER`         | `nexus`                 | PostgreSQL user           |
| `NEXUS_DB_PASSWORD`     | *(required)*            | PostgreSQL password       |
| `NEXUS_DB_NAME`         | `nexus`                 | PostgreSQL database       |
| `NEXUS_DB_SSLMODE`      | `disable`               | PostgreSQL SSL mode       |
| `NEXUS_JWT_SECRET`      | `change-me-in-production` | HMAC secret for JWT    |
| `NEXUS_LOG_LEVEL`       | `info`                  | debug, info, warn, error  |
| `NEXUS_LOG_FORMAT`      | `json`                  | json or text              |

## Project Structure

```
cmd/server/      — Server bootstrap, graceful shutdown
internal/
  auth/          — JWT generation, validation, token hashing
  config/        — Environment-based configuration
  db/            — Database layer (planned)
  handler/       — HTTP request handlers
  middleware/     — Logging, CORS, auth, recovery
  models/        — Domain types (User, Session, ConnectedApp, etc.)
```

## Development

```bash
# Run with Go
go run .

# Run with live reload
make dev

# Tests
make test

# Lint
make lint
```

