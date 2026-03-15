# Architecture & Core Rules: project_abc (MVP)

## Objective
Build a micro-SaaS booking and SMS notification system for the beauty industry. Focus on minimizing "no-shows" via automated reminders and "Magic Links" for appointment confirmation/cancellation.

## Core Directives (DevOps First)
1. **Everything is an Environment Variable:** Never hardcode project names, ports, DB credentials, or API keys. Use `.env`. The project name should be dynamically loaded as `PROJECT_NAME`.
2. **Tech Stack:**
    - Backend: Golang (v1.22+).
    - Frontend: SvelteKit + Tailwind CSS (v3/v4).
    - DB: PostgreSQL.
    - Infra: Docker + Docker Compose.
3. **Monorepo Structure:**
    - `/backend` (Go API)
    - `/frontend` (SvelteKit)
    - `docker-compose.yml` at the root level.
4. **Prod-Ready Standards:** Use structured JSON logging (`log/slog` in Go). Explicit error handling.