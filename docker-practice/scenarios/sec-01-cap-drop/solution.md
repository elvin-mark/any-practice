# Solution: Least Privilege - Linux Capability Dropping

### Key Concepts:
- Traditional POSIX root divides privileges into individual Linux Capabilities (e.g., `CAP_NET_BIND_SERVICE`, `CAP_SYS_PTRACE`, `CAP_NET_RAW`).
- The most secure defensive posture is to drop **ALL** capabilities, then selectively white-list only the specific ones required for application execution.

### Command:
```bash
docker run -d \
  --name least-priv-web \
  --cap-drop=ALL \
  --cap-add=NET_BIND_SERVICE \
  --cap-add=CHOWN \
  --cap-add=SETUID \
  --cap-add=SETGID \
  -p 8089:80 \
  --label practice=docker \
  nginx:alpine
```
