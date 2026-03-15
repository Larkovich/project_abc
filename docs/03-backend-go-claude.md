# Backend (Golang) Requirements

## Rules
1. Load environment variables using a package like `godotenv` or read directly from `os.Getenv` after Docker injects them.
2. Standard layout:
    - `/cmd/api/main.go` (Entrypoint).
    - `/internal/api` (HTTP Handlers / Chi router).
    - `/internal/database` (DB connection pooling using `database/sql`, avoid heavy ORMs).
    - `/internal/worker` (Background goroutine for checking `notifications` table).
3. Start with an MVP: Add a `GET /health` endpoint that returns `{"status": "ok", "project": "<PROJECT_NAME>"}`.