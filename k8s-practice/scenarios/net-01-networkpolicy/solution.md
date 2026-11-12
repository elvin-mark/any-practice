# Solution: Restrict Ingress Traffic with NetworkPolicy

### Step 1: Create NetworkPolicy Manifest
Save to `db-policy.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: db-access-policy
  namespace: secure-zone
spec:
  podSelector:
    matchLabels:
      role: db
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: backend
    ports:
    - protocol: TCP
      port: 80
```

### Step 2: Apply the Manifest
```bash
kubectl apply -f db-policy.yaml
```

### Step 3: Test Connectivity
```bash
DB_IP=$(kubectl get pod db -n secure-zone -o jsonpath='{.status.podIP}')

# Test from backend (allowed):
kubectl exec backend -n secure-zone -- wget -qO- -T 3 "http://$DB_IP"

# Test from frontend (blocked / times out):
kubectl exec frontend -n secure-zone -- wget -qO- -T 3 "http://$DB_IP"
```
