# Solution: Multi-Tier Compose Stack with Health Dependencies

### Key Concepts:
- In Compose v2, `depends_on` supports health condition checking:
  ```yaml
  depends_on:
    redis:
      condition: service_healthy
  ```
- This ensures the `web` container is only started AFTER the Redis health check returns exit code 0 (`PONG`).

### Reference `docker-compose.yml`:
```yaml
services:
  redis:
    image: redis:alpine
    networks:
      - stack-net
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3

  web:
    image: nginx:alpine
    ports:
      - "8088:80"
    networks:
      - stack-net
    depends_on:
      redis:
        condition: service_healthy

networks:
  stack-net:
    driver: bridge
```

### Launch Command:
```bash
docker compose up -d
```
