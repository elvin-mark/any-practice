# Solution: Configure RBAC for Developer Role

### Step 1: Create the Role
```bash
kubectl create role developer-role \
  --namespace=development \
  --verb=get,list,watch,create,update,delete \
  --resource=pods,services
```

### Step 2: Create the RoleBinding
```bash
kubectl create rolebinding developer-binding \
  --namespace=development \
  --role=developer-role \
  --user=developer
```

### Step 3: Test and Verify Permissions
```bash
# Verify user developer can list pods in development
kubectl auth can-i list pods --namespace=development --as=developer
# Expected output: yes

# Verify user developer can create services in development
kubectl auth can-i create services --namespace=development --as=developer
# Expected output: yes

# Verify user developer cannot access kube-system
kubectl auth can-i list pods --namespace=kube-system --as=developer
# Expected output: no
```
