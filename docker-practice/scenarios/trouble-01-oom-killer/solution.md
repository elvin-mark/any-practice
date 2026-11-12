# Solution: Diagnose and Fix OOMKilled Container

### Key Concepts:
- In Linux, when a cgroup exceeds its memory quota and cannot reclaim free pages, the kernel invokes the OOM killer (`oom-killer`) to terminate the process via `SIGKILL` (signal 9).
- Docker flags this in container state: `ExitCode: 137` (128 + 9), `OOMKilled: true`.
- Diagnosing OOM:
  ```bash
  docker inspect mem-hungry --format 'Status={{.State.Status}} ExitCode={{.State.ExitCode}} OOM={{.State.OOMKilled}}'
  ```

### Resolution Commands:
```bash
# 1. Remove crashed container
docker rm -f mem-hungry

# 2. Re-run with adequate memory limit (256MB)
docker run -d \
  --name mem-hungry \
  --memory=256m \
  --memory-swap=256m \
  --label practice=docker \
  alpine:3.20 sh -c 'dd if=/dev/zero of=/dev/null bs=1M count=100 && sleep 3600'
```
