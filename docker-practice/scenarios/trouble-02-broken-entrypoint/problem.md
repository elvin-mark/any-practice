# Task: Fix Exec Format & Missing Shebang Entrypoint

A newly written image failed on execution with:
`exec /entrypoint.sh: no such file or directory` or `permission denied` (Exit Code 127).

### Objectives:
1. Examine `entrypoint.sh` and `Dockerfile` in this directory.
2. Troubleshoot and resolve the entrypoint issues:
   - Ensure `entrypoint.sh` includes a valid shebang line (`#!/bin/sh`).
   - Ensure `entrypoint.sh` has executable permissions (`chmod +x entrypoint.sh`).
   - Verify Unix LF line endings.
3. Build the corrected image: `docker build -t fixed-entrypoint:latest .`
4. Run the container: `docker run -d --name fixed-worker --label practice=docker fixed-entrypoint:latest`
5. The container must start cleanly without exiting with code 127, and its logs must contain `ENTRYPOINT_EXECUTION_SUCCESSFUL`.
