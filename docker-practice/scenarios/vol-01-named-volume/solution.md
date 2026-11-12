# Solution: Persistent Storage with Named Volumes

### Key Concepts:
- Named volumes are completely managed by Docker and stored in `/var/lib/docker/volumes/` on Linux hosts.
- They are easier to back up, migrate, and share among containers than bind mounts.
- Even if the container is destroyed (`docker rm -f datastore-app`), the data in `db-store` remains intact until explicitly deleted.

### Commands:
```bash
# 1. Create the named volume
docker volume create --label practice=docker db-store

# 2. Run container with volume mounted
docker run -d \
  --name datastore-app \
  -v db-store:/var/lib/appdata \
  --label practice=docker \
  alpine:3.20 sleep 3600

# 3. Create persistent data file
docker exec datastore-app sh -c 'echo DOCKER_PERSISTENCE_VERIFIED > /var/lib/appdata/persistent.txt'
```
