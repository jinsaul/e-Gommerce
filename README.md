# e-Gommerce

A mock e-commerce website built with Go, Angular, and MongoDB.

## Stack

- **Backend:** Go (stdlib `net/http`)
- **Frontend:** Angular 19 (Signals, standalone components)
- **Database:** MongoDB Atlas
- **Containerization:** Docker & Docker Compose

## Run

```bash
# Docker
docker compose up --build

# Or locally
cd server && go run .     # API on :8080
cd client && ng serve     # App on :4200
```
