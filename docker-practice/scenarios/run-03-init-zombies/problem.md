# Task: Fix Zombie Process Reaping with Docker Init

When a containerized process spawns child processes that exit without their parent waiting on them, they become zombie (`defunct`) processes. Under traditional Linux, PID 1 (`init` / `systemd`) reaps them. In Docker, your application runs as PID 1 by default, leading to process table exhaustion and failed signal handling.

### Objectives:
1. Launch a container named `init-worker` using `alpine:3.20`.
2. Configure Docker's built-in init system using the `--init` flag so that `docker-init` (tini) runs as PID 1 inside the container.
3. Keep the container running in the background with `sleep 3600`.
4. Attach label `practice=docker`.
5. Verify via `docker inspect` that `HostConfig.Init` is set to `true`, and check inside the container that PID 1 is `/sbin/docker-init` (or `docker-init`).
