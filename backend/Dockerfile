# ---- Dev stage: hot-reload with go run ----
FROM golang:1.24-alpine AS dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "cmd/api/main.go"]

# ---- Builder stage: compile static binary ----
FROM dev AS builder
RUN go build -o api cmd/api/main.go

# ---- Prod stage: minimal runtime image ----
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]
