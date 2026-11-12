# Solution: User-Defined Bridge Network & DNS Discovery

### Key Concepts:
- The default bridge network (`bridge`) lacks automatic name resolution. Docker's embedded DNS server (`127.0.0.11`) is only activated on user-defined networks.
- Specifying `--subnet` guarantees predictable IP addressing across container environments.

### Commands:
```bash
# 1. Create network with specified subnet
docker network create \
  --subnet 172.28.0.0/16 \
  --label practice=docker \
  custom-bridge

# 2. Run API server
docker run -d \
  --name api-server \
  --network custom-bridge \
  --label practice=docker \
  nginx:alpine

# 3. Run client application
docker run -d \
  --name client-app \
  --network custom-bridge \
  --label practice=docker \
  alpine:3.20 sleep 3600
```
