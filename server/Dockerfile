# ---- Build stage ----
# Uses the official Go image to compile the app
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first (Docker caches this layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o server .

# ---- Run stage ----
# Uses a tiny image (just Linux, no Go) to run the compiled binary
FROM alpine:3.21

WORKDIR /app

# Copy only the compiled binary from the build stage
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
