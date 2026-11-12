# Task: Least Privilege - Linux Capability Dropping

Containers are essentially isolated Linux processes running within namespaces and cgroups. Granting unnecessary kernel capabilities (e.g. `CAP_SYS_ADMIN`, `CAP_NET_RAW`, `CAP_MKNOD`) significantly increases the risk of container breakout and kernel privilege escalation.

### Objectives:
Run a hardened container named `least-priv-web` using `nginx:alpine` following the Least Privilege model:
1. **Container Name**: `least-priv-web`
2. **Base Image**: `nginx:alpine`
3. **Capability Drop**: Drop **ALL** default Linux capabilities (`--cap-drop=ALL`).
4. **Capability Add**: Add back **ONLY**:
   - `NET_BIND_SERVICE` (to bind to privileged network ports < 1024)
   - `CHOWN` (for nginx file ownership setup)
   - `SETUID` and `SETGID` (for worker process privilege separation)
5. **Port Binding**: Map `-p 8089:80`.
6. **Label**: `practice=docker`
7. Ensure `least-priv-web` is running and responds with HTTP 200 on `http://127.0.0.1:8089`.
