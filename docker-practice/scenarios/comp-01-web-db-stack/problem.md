# Task: Multi-Tier Compose Stack with Health Dependencies

A common microservice failure in Docker Compose occurs when a web frontend starts up before its dependent database or cache is fully initialized and accepting connections, causing crash loops.

### Objectives:
1. In this directory, construct a valid `docker-compose.yml` defining two services:
   - **`redis`**:
     - Image: `redis:alpine`
     - Healthcheck:
       - Test: `["CMD", "redis-cli", "ping"]`
       - Interval: `5s`
       - Timeout: `3s`
       - Retries: `3`
     - Attached to custom network `stack-net`.
   - **`web`**:
     - Image: `nginx:alpine`
     - Depends on `redis` with condition `service_healthy`.
     - Port mapping: `8088:80`.
     - Attached to custom network `stack-net`.
2. Launch the stack in detached mode: `docker compose up -d`.
3. Verify that:
   - Both services are up and running.
   - `redis` status is `healthy`.
   - `http://127.0.0.1:8088` returns HTTP 200.
