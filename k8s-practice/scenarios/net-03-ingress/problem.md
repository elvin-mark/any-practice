# Task: Configure Ingress with Path-Based Routing

The `ecommerce` namespace contains two microservices:
- `catalog-svc` (port 80)
- `checkout-svc` (port 80)

Create an Ingress resource to route HTTP traffic:
1. Ingress Name: `ecommerce-ingress`
2. Namespace: `ecommerce`
3. Host: `shop.example.com`
4. Path Rules:
   - Path `/catalog` with pathType `Prefix` must forward to `catalog-svc` on port `80`.
   - Path `/checkout` with pathType `Prefix` must forward to `checkout-svc` on port `80`.
5. Ensure the ingress resource applies cleanly with valid `networking.k8s.io/v1` syntax.
