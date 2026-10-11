# System Design Implementations

A collection of hands-on, educational implementations of common distributed system design patterns and components, written in Go and Python.

> [!NOTE]
> Created for educational and learning purposes to understand low-level architecture, concurrency, and distributed algorithms.

---

## Modules

### 1. ⏱️ Rate Limiter Gateway (`ratelimiter/`)
A distributed API Rate Limiter and reverse proxy written in Go:
- **Algorithm**: Token Bucket algorithm (`internal/limiter/tokenbucket`) for burst management and steady-state rate control.
- **Storage**: Redis-backed state store (`internal/store/redis`) for distributed limit tracking across cluster instances.
- **Middleware & Proxy**: Reverse proxy integration (`internal/proxy`) with HTTP middleware for request enrichment, logging, and metric export.
- **Documentation**: Integrated Swagger API specification (`docs/`).

### 2. 📰 Twitter Feed System (`twitter-feed/`)
An architectural implementation of a social feed generation and distribution pipeline:
- **Feed Gateway (`twitter-feed/gateway/`)**: FastAPI gateway for posting tweets and retrieving timeline feeds.
- **Containerization**: Docker Compose definition for running local backing services.

---

## Project Structure

```
system-design-implementation/
├── ratelimiter/                # Go distributed token bucket rate limiter
│   ├── cmd/gateway/            # Gateway entrypoint
│   ├── internal/limiter/       # Token bucket algorithm
│   ├── internal/store/         # Redis store integration
│   └── docs/                   # Swagger documentation
└── twitter-feed/               # Social feed architecture
    └── gateway/                # FastAPI feed gateway service
```

---

## Getting Started

### Rate Limiter (Go)

```bash
cd ratelimiter
go mod download
go run cmd/gateway/main.go
```

### Twitter Feed (Python)

```bash
cd twitter-feed/gateway
docker compose up -d
uv run uvicorn app.main:app --reload
```

---

## License

Educational / Personal Study.
