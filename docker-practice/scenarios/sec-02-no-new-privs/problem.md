# Task: Prevent SUID Escalation with No-New-Privileges

Setuid binaries (such as `passwd`, `sudo`, or malicious suid payloads) allow an unprivileged user process to elevate privileges to the binary's owner (often root) upon execution. Setting the `no-new-privileges` flag disables execution of setuid and setgid bits completely.

### Objectives:
1. Run a container named `no-escalate` using `alpine:3.20 sleep 3600` in the background.
2. Enable the `no-new-privileges` security flag via `--security-opt=no-new-privileges:true`.
3. Attach label `practice=docker`.
4. Verify using `docker inspect` that `SecurityOpt` contains `no-new-privileges:true` (or `no-new-privileges`).
