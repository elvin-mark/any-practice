# Solution: Multi-Stage Build for Go Microservice

### Key Concepts:
1. **Multi-Stage Builds**: Allow dividing the Docker build process into discrete stages. The build environment (SDKs, build tools, package managers) is discarded, and only the final compiled artifact is copied to a minimal runtime image.
2. **Static Compilation**: Go binaries can be built without C standard library dependencies (`CGO_ENABLED=0`), allowing them to run on minimal base images such as Alpine or even `scratch`.
3. **Security (Non-Root User)**: Containers should never run as root unless strictly necessary. Creating a dedicated UID (e.g. `10001`) prevents container escape vulnerabilities.

### Reference `Dockerfile`:
```dockerfile
# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY main.go ./

# Compile static binary with debugging symbols stripped (-s -w)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# Stage 2: Minimal runtime stage
FROM alpine:3.20

# Create unprivileged user
RUN adduser -D -u 10001 -s /sbin/nologin appuser

WORKDIR /app
COPY --from=builder /app/server /app/server

USER 10001:10001
EXPOSE 8080

ENTRYPOINT ["/app/server"]
```

### Build Command:
```bash
docker build -t go-microservice:optimized .
```
