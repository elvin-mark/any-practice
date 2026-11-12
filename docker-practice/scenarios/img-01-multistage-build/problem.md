# Task: Multi-Stage Build for Minimal Go Microservice

A team member deployed a Go microservice using a single-stage `Dockerfile.orig` that includes the full Go toolchain. As a result, the image size exceeds 800 MB and runs as the root user.

### Objectives:
1. Examine the Go HTTP server in `main.go`.
2. Write an optimized multi-stage `Dockerfile` in this directory:
   - **Stage 1 (Builder)**: Compile the application into a statically linked binary using `CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .`.
   - **Stage 2 (Runtime)**: Copy ONLY the compiled binary into a minimal runtime base (`alpine:3.20` or `scratch`).
   - Create an unprivileged user (UID `10001`, name `appuser`) and configure `USER 10001`.
   - Expose port `8080` and set `ENTRYPOINT ["/app/server"]`.
3. Build and tag the resulting image as `go-microservice:optimized`.
4. The final image size must be strictly less than **35 MB**.
5. When run, the container must respond with HTTP 200 on port 8080 and run as an unprivileged user.
