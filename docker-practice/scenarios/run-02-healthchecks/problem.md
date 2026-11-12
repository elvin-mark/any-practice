# Task: Implement Container Healthchecks

Docker container status usually only reflects whether process PID 1 is alive, even if the underlying web server or database has hung or is returning 500 errors. Healthchecks let the Docker daemon monitor actual service health.

### Objectives:
Run a detached container named `monitored-web` using `nginx:alpine` with the following healthcheck parameters:
1. **Container Name**: `monitored-web`
2. **Image**: `nginx:alpine`
3. **Healthcheck Command**: `wget --no-verbose --tries=1 --spider http://127.0.0.1/ || exit 1`
4. **Interval**: `5s` (`--health-interval=5s`)
5. **Timeout**: `3s` (`--health-timeout=3s`)
6. **Retries**: `3` (`--health-retries=3`)
7. **Start Period**: `2s` (`--health-start-period=2s`)
8. **Label**: `practice=docker`

When properly configured, `docker inspect monitored-web --format '{{.State.Health.Status}}'` must transition to `healthy`.
