[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/yIQ32uio)

# p2-final-project-backend-team

**Short description**: Backend API for the final project (Go + Echo + GORM + Postgres), with Docker support.

---

## 📁 Project Structure
High-level overview of the main folders and their purpose:

- `cmd/` — entry points
  - `api/` — API server
  - `migration/` — database migration runner
- `internal/` — application code
  - `configs/` — configuration loader (`NewConfig()` loads `.env` / env vars)
  - `domain/` — domain models (e.g., `models.User`)
  - `dto/` — request/response DTOs
  - `handlers/` — HTTP handlers
  - `repositories/` — data access layer (GORM)
  - `services/` — business logic
  - `infrastructure/` — DB connection, context helpers
  - `router/` — route registration endpoint
- `docker/` — Dockerfiles and build assets
- `docker-compose.yml` — local compose for API + Postgres
- `Makefile` — convenient tasks (`start-app`, `migrate-schema`, `wire-gen`)

---

## ⚙️ Setup & Installation
Quick setup steps to run the project locally.

### Prerequisites
- Docker & Docker Compose
- Go (optional, for running migrations with `go run`)
- (Optional) GNU Make

### 1) Create `.env` (copy from `.env.example`)
You can create `.env` by copying the provided example file and then editing it as needed.

```bash
cp .env.example .env
```

Then open `.env` and fill in the required values. Example contents:

```env
PORT=8080
DATABASE_HOST=postgres
DATABASE_NAME=dev_db
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=secret
DATABASE_PORT=5432
```

Note: `configs.NewConfig()` uses `godotenv` to load `.env` and parse values.

### 2) Start services with Docker Compose
```bash
make start-app
# or
docker-compose up --build -d
```
This builds the API image and starts the API + Postgres containers.

### 3) Run database migrations
Ensure the DB is reachable, then run:

```bash
make migrate-schema
# or

go run ./cmd/migration/main.go
```

Note: If you want to run migrations from inside the container, set `DATABASE_HOST` to the Postgres service name (for example `postgres`) or update `cmd/migration/main.go` to read the host from env.


## 🚀 Quick API examples
Base URL: `http://localhost:<PORT>/api/v1`

- Create user (POST `/users`)
```bash
curl -i -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{"nisit_id":"63060001","user_fullname":"Somchai S.","email":"somchai@example.com","role":"student"}'
```

- Get user by NisitID (GET `/users/:nisit_id`)
```bash
curl -i "http://localhost:8080/api/v1/users/63060001"
```


---

