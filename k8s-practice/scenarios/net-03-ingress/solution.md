# Solution: Configure Ingress with Path-Based Routing

### Step 1: Create the Ingress YAML
Save to `ecommerce-ingress.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ecommerce-ingress
  namespace: ecommerce
spec:
  rules:
  - host: shop.example.com
    http:
      paths:
      - path: /catalog
        pathType: Prefix
        backend:
          service:
            name: catalog-svc
            port:
              number: 80
      - path: /checkout
        pathType: Prefix
        backend:
          service:
            name: checkout-svc
            port:
              number: 80
```

### Step 2: Apply the Ingress
```bash
kubectl apply -f ecommerce-ingress.yaml
```

### Step 3: Verify with kubectl
```bash
kubectl get ingress ecommerce-ingress -n ecommerce
kubectl describe ingress ecommerce-ingress -n ecommerce
```
