# Solution: Secure Host IP Port Binding

### Key Concepts:
- The standard syntax `-p [host_ip:][host_port:]container_port` determines the interface binding.
- Omitting `host_ip` binds to `0.0.0.0` (all addresses), allowing anyone on the same network or public internet to access the port.
- Explicitly binding to `127.0.0.1` restricts network traffic strictly to local host processes.

### Command:
```bash
docker run -d \
  --name internal-cache \
  -p 127.0.0.1:6379:6379 \
  --label practice=docker \
  redis:alpine
```
