# Task: Restrict Ingress Traffic with NetworkPolicy

Namespace `secure-zone` hosts three components:
- `db` (listening on TCP port 80, labeled `role=db`)
- `backend` (labeled `role=backend`)
- `frontend` (labeled `role=frontend`)

Configure network isolation:
1. Create a `NetworkPolicy` named `db-access-policy` in namespace `secure-zone`.
2. The policy must isolate incoming traffic to pods labeled `role=db`.
3. Allow ingress traffic to `db` **strictly** from pods labeled `role=backend` on TCP port `80`.
4. Block all other ingress traffic to `db` (e.g. connections originating from `frontend`).
