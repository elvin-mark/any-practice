# Solution: Hardened Read-Only Rootfs with Writable Tmpfs

### Key Concepts:
- `--read-only`: Mounts the container rootfs as read-only. Attempts to write or modify any files outside explicit volume mounts fail with `Read-only file system`.
- `--tmpfs <path>`: Mounts an ephemeral Linux `tmpfs` (in RAM) at the designated path. Files written here are fast, never written to disk, and disappear when the container stops.
- This pattern satisfies CIS Docker Benchmark 5.12: "Mount container's root filesystem as read-only".

### Command:
```bash
docker run -d \
  --name secure-nginx \
  --read-only \
  --tmpfs /var/run \
  --tmpfs /var/cache/nginx \
  --tmpfs /tmp \
  -p 8085:80 \
  --label practice=docker \
  nginx:alpine
```
