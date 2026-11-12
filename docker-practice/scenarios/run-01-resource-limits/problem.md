# Task: Enforce CPU, Memory & Restart Policies

In production, runaway containers without resource constraints can exhaust host memory and CPU, triggering kernel OOM panics.

### Objectives:
Run a detached container named `prod-worker` with the following strict resource and policy constraints:
1. **Container Name**: `prod-worker`
2. **Base Image & Command**: `alpine:3.20 sleep 3600`
3. **Memory Limit**: Exactly `256m` (268,435,456 bytes)
4. **Memory Swap Limit**: Exactly `256m` (to disable swap extension beyond the memory limit)
5. **CPU Quota**: Exactly `0.5` CPUs (equivalent to `NanoCpus: 500000000`)
6. **PID Limit**: Maximum of `100` processes (`--pids-limit=100`)
7. **Restart Policy**: `unless-stopped`
8. **Label**: `practice=docker`
