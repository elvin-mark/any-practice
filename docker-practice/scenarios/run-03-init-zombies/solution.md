# Solution: Fix Zombie Process Reaping with Docker Init

### Key Concepts:
- The Linux kernel assigns special status to Process ID 1. When child processes terminate, they stay in the process table in a zombie (`Z`) state until their parent calls `wait()` or `waitpid()`.
- If the parent dies first, orphaned children are reparented to PID 1. If PID 1 does not reap them, they accumulate and leak kernel PID table entries.
- The `--init` flag inserts `tini` as PID 1, which reaps zombie processes and correctly forwards UNIX signals (like `SIGTERM`, `SIGINT`) to child processes.

### Command:
```bash
docker run -d \
  --name init-worker \
  --init \
  --label practice=docker \
  alpine:3.20 sleep 3600
```
