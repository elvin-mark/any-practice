# Solution: Prevent SUID Escalation with No-New-Privileges

### Key Concepts:
- The Linux kernel flag `PR_SET_NO_NEW_PRIVS` prevents a process and its descendants from gaining additional privileges through `execve()` calls.
- Even if an attacker finds an existing setuid binary or places a malicious setuid binary into a writable volume, execution will NOT escalate privileges.

### Command:
```bash
docker run -d \
  --name no-escalate \
  --security-opt=no-new-privileges:true \
  --label practice=docker \
  alpine:3.20 sleep 3600
```
