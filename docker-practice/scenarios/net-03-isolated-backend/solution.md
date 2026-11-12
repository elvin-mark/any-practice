# Solution: Multi-Tier Network Isolation (DMZ vs Backend)

### Key Concepts:
- Docker user-defined networks create isolated Linux bridge interfaces (`br-<network-id>`).
- Containers connected to different bridges cannot exchange packets unless connected to the same bridge network or explicitly routed.
- Connecting an application gateway to both networks allows it to route traffic securely between clients and backend databases while maintaining network isolation.

### Commands:
```bash
# 1. Create the two isolated networks
docker network create --label practice=docker frontend-net
docker network create --label practice=docker backend-net

# 2. Run backend database on backend-net only
docker run -d \
  --name secure-db \
  --network backend-net \
  --label practice=docker \
  alpine:3.20 sleep 3600

# 3. Run external client on frontend-net only
docker run -d \
  --name external-client \
  --network frontend-net \
  --label practice=docker \
  alpine:3.20 sleep 3600

# 4. Run gateway on frontend-net and attach to backend-net
docker run -d \
  --name app-gateway \
  --network frontend-net \
  --label practice=docker \
  alpine:3.20 sleep 3600

docker network connect backend-net app-gateway
```
