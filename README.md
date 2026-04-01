# e-Gommerce

A mock e-commerce website built with Go, Angular, and MongoDB.

## Stack

- **Backend:** Go (stdlib `net/http`)
- **Frontend:** Angular 19 (Signals, standalone components)
- **Database:** MongoDB Atlas
- **Containerization:** Docker & Docker Compose

## Prerequisites

- [Go 1.22+](https://go.dev/doc/install)
- [Node.js & npm](https://nodejs.org/en/download)
- MongoDB Atlas cluster (add your IP to the allowlist, copy connection string into `.env`)

## Run

```bash
# Docker
docker compose up --build

# Or locally
cd server && go run .     # API on :8080
cd client && ng serve     # App on :4200
```

