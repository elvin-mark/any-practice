# Solution: Optimize Dockerfile Layer Caching

### Key Concepts:
Docker creates intermediate image layers for each instruction (`COPY`, `RUN`, `ADD`).
If an instruction's inputs haven't changed, Docker reuses the cached layer from previous builds.
Because source code (`server.js`) changes far more frequently than external dependencies (`package.json`), copying the manifest file first allows Docker to cache the expensive installation layer across subsequent source code changes.

### Reference `Dockerfile`:
```dockerfile
FROM alpine:3.20

WORKDIR /app

# 1. Copy package files first
COPY package*.json ./

# 2. Run dependency installation (this layer will now be cached)
RUN echo "Installing deps..." && touch /app/deps_installed

# 3. Copy application source code
COPY . .

CMD ["cat", "/app/deps_installed"]
```

### Build Command:
```bash
docker build -t cached-app:latest .
```
