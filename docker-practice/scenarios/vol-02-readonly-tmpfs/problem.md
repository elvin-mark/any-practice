# Task: Hardened Read-Only Rootfs with Writable Tmpfs

Running containers with a read-only root filesystem prevents attackers or malicious payloads from modifying binary files, installing rootkits, or creating unauthorized persistent backdoors. However, many server applications (like Nginx) still require volatile temporary storage to write runtime PIDs and caches.

### Objectives:
Run a hardened container named `secure-nginx` using `nginx:alpine` satisfying:
1. **Container Name**: `secure-nginx`
2. **Image**: `nginx:alpine`
3. **Read-Only Root Filesystem**: `--read-only` flag enabled.
4. **Tmpfs Mounts**: Provide writable in-memory `tmpfs` mounts for:
   - `/var/run`
   - `/var/cache/nginx`
   - `/tmp`
5. **Port Binding**: Map port `8085` on host to port `80` in container (`-p 8085:80`).
6. **Label**: `practice=docker`
7. When verified, `secure-nginx` must be running, have `ReadonlyRootfs == true`, and respond with HTTP 200 on `http://127.0.0.1:8085`.
