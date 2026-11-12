# Task: Multi-Tier Network Isolation (DMZ vs Backend)

In secure multi-tier architectures, backend data stores must be completely isolated from external networks and ingress clients. Only an intermediary gateway / application proxy should be capable of communicating with both frontend clients and backend storage.

### Objectives:
1. Create two user-defined bridge networks labeled `practice=docker`:
   - `frontend-net`
   - `backend-net`
2. Run a container named `secure-db` using `alpine:3.20 sleep 3600` connected **ONLY** to `backend-net`.
3. Run a container named `external-client` using `alpine:3.20 sleep 3600` connected **ONLY** to `frontend-net`.
4. Run a container named `app-gateway` using `alpine:3.20 sleep 3600` connected to **BOTH** `frontend-net` and `backend-net`.
5. Attach label `practice=docker` to all created resources.
6. Verification tests:
   - `app-gateway` can ping both `external-client` and `secure-db` by container name.
   - `external-client` **CANNOT** reach `secure-db` (network isolation verified).
