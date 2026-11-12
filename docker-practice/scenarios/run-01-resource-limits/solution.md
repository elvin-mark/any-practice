# Solution: Enforce CPU, Memory & Restart Policies

### Key Concepts:
- `--memory=256m`: Sets maximum RAM available to the container cgroup.
- `--memory-swap=256m`: Setting `--memory-swap` equal to `--memory` disables additional swap allocation, ensuring strict memory ceiling.
- `--cpus=0.5`: Guarantees the container can use at most 50% of a single CPU core.
- `--pids-limit=100`: Mitigates fork-bomb attacks by limiting concurrent processes.
- `--restart=unless-stopped`: Container restarts on failure or Docker daemon reboot, unless manually stopped.

### Command:
```bash
docker run -d \
  --name prod-worker \
  --memory=256m \
  --memory-swap=256m \
  --cpus=0.5 \
  --pids-limit=100 \
  --restart=unless-stopped \
  --label practice=docker \
  alpine:3.20 sleep 3600
```
