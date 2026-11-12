# Task: Harden Dockerfile Non-Root User & Permissions

The production audit found that the application image runs under the `root` user, which violates the Principle of Least Privilege and organizational security policy.

### Objectives:
1. Review the existing `Dockerfile` in this directory.
2. Update `Dockerfile` to enforce non-root execution:
   - Create a dedicated non-root user `secuser` with UID `10001` and GID `10001`.
   - Set ownership of the `/app` and `/data` directories to `secuser:secuser`.
   - Ensure `secuser` has write permissions to `/data`.
   - Switch active user using `USER 10001`.
3. Build and tag the image as `hardened-app:latest`.
4. When run, the container must verify that:
   - `id -u` outputs `10001` (not `0`).
   - The application can create and write files inside `/data`.
