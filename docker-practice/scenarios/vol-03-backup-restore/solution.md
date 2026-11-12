# Solution: Volume Backup & Restoration Pipeline

### Key Concepts:
- Docker volumes are stored inside host paths protected by Docker daemon permissions.
- The standard, robust mechanism for backing up and restoring volumes without stopping Docker is using ephemeral containers with volume mounts.

### Commands:
```bash
# 1. Backup: Mount source volume and host directory, create compressed tarball
docker run --rm \
  -v prod-backup-source:/volume:ro \
  -v /tmp:/backup \
  alpine:3.20 \
  tar czf /backup/prod-backup.tar.gz -C /volume .

# 2. Create restoration target volume
docker volume create --label practice=docker prod-backup-restored

# 3. Restore: Mount target volume and host archive directory, extract tarball
docker run --rm \
  -v prod-backup-restored:/volume \
  -v /tmp:/backup \
  alpine:3.20 \
  tar xzf /backup/prod-backup.tar.gz -C /volume
```
