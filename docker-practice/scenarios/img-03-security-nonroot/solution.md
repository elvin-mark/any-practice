# Solution: Harden Dockerfile Non-Root User & Permissions

### Key Concepts:
- Running as `root` in a container means that a container breakout may grant root privileges directly on the host kernel.
- In Linux, unprivileged application users are assigned UIDs >= 1000 or >= 10000.
- Creating the user, setting ownership with `chown -R`, and declaring `USER <uid>` guarantees that all runtime commands execute with restricted permissions.

### Reference `Dockerfile`:
```dockerfile
FROM alpine:3.20

WORKDIR /app

# Create unprivileged group and user (UID 10001)
RUN addgroup -g 10001 secuser && \
    adduser -D -u 10001 -G secuser -s /bin/sh secuser && \
    mkdir -p /app /data && \
    chown -R secuser:secuser /app /data

COPY app.sh /app/app.sh
RUN chmod +x /app/app.sh

# Switch to the non-root user
USER 10001

CMD ["/app/app.sh"]
```

### Build Command:
```bash
docker build -t hardened-app:latest .
```
