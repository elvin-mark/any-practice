# Task: User-Defined Bridge Network & DNS Discovery

On Docker's default `bridge` network, containers can only communicate via direct IP addresses (unless legacy `--link` is used). User-defined bridge networks provide automatic embedded DNS resolution so containers can address each other by container name.

### Objectives:
1. Create a user-defined bridge network named `custom-bridge` with subnet `172.28.0.0/16` and label `practice=docker`.
2. Run a web server container named `api-server` attached to `custom-bridge` using `nginx:alpine` and label `practice=docker`.
3. Run a client container named `client-app` attached to `custom-bridge` using `alpine:3.20` in the background (`sleep 3600`) and label `practice=docker`.
4. Ensure that `client-app` can resolve and connect to `http://api-server` using Docker's embedded DNS server (`127.0.0.11`).
