# Task: Diagnose and Fix OOMKilled Container

A critical background analytics worker named `mem-hungry` crashed immediately after startup. The operations team reported exit code 137.

### Objectives:
1. Inspect the stopped container `mem-hungry` using `docker inspect` to confirm whether `State.OOMKilled` is `true`.
2. Delete or recreate `mem-hungry` with a sufficient memory allocation:
   - Base Image & Command: `alpine:3.20 sh -c 'dd if=/dev/zero of=/dev/null bs=1M count=100 && sleep 3600'`
   - Set memory limit to `256m` (`--memory=256m --memory-swap=256m`).
   - Run in detached mode (`-d`).
   - Label: `practice=docker`.
3. Verify that `mem-hungry` remains running without being terminated by the host Out-Of-Memory (OOM) killer.
