# Solution: Implement Container Healthchecks

### Key Concepts:
- Docker healthchecks verify application layer readiness instead of merely checking OS process termination.
- Container states transition from `starting` -> `healthy` (or `unhealthy` if retries fail).
- In Docker Compose and Swarm, orchestrators use health status for zero-downtime rolling deploys and `depends_on` conditions.

### Command:
```bash
docker run -d \
  --name monitored-web \
  --health-cmd="wget --no-verbose --tries=1 --spider http://127.0.0.1/ || exit 1" \
  --health-interval=5s \
  --health-timeout=3s \
  --health-retries=3 \
  --health-start-period=2s \
  --label practice=docker \
  nginx:alpine
```
