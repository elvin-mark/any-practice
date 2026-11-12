# Solution: Fix Exec Format & Missing Shebang Entrypoint

### Key Concepts:
- In Linux, an executable script requires:
  1. A valid shebang (`#!/bin/sh`) on line 1 informing the kernel which interpreter to invoke.
  2. Executable permission bit (`chmod +x`).
  3. Unix line endings (`LF`, not DOS `CRLF` `\r\n`). When a script has Windows CRLF, the Linux kernel looks for `/bin/sh\r`, resulting in `no such file or directory`.

### Resolution:
1. Update `entrypoint.sh`:
```bash
#!/bin/sh
echo 'ENTRYPOINT_EXECUTION_SUCCESSFUL'
sleep 3600
```

2. Make executable:
```bash
chmod +x entrypoint.sh
```

3. Build and run:
```bash
docker build -t fixed-entrypoint:latest .
docker run -d --name fixed-worker --label practice=docker fixed-entrypoint:latest
```
