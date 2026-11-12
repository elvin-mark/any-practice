# Task: Secure Host IP Port Binding

When publishing ports with `docker run -p 6379:6379`, Docker automatically binds the port to `0.0.0.0` (all host network interfaces), exposing internal caches and databases to the public Internet or WAN without authentication.

### Objectives:
1. Run a container named `internal-cache` using `redis:alpine` in detached mode.
2. Bind the container port `6379` strictly to the host's loopback interface `127.0.0.1` on host port `6379` (`-p 127.0.0.1:6379:6379`).
3. Attach label `practice=docker`.
4. Verify using `docker inspect` that `HostIp` is `127.0.0.1` and `HostPort` is `6379`.
5. Verify that Redis responds to a ping on `127.0.0.1:6379`.
