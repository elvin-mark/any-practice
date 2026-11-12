# Solution: Create a Static Pod on Control Plane Node

### Step 1: Identify the Control-Plane Node
```bash
kubectl get nodes
# Identify node name, e.g. cka-cluster-1-control-plane
```

### Step 2: Access the Node Container
```bash
docker exec -it <control-plane-node-name> bash
# or: docker exec -it $(docker ps --filter "name=control-plane" -q | head -n 1) bash
```

### Step 3: Generate and Place the Static Pod Manifest
Inside the control-plane container:
```bash
cat <<EOF > /etc/kubernetes/manifests/my-static-web.yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-static-web
  namespace: default
spec:
  restartPolicy: Always
  containers:
  - name: my-static-web
    image: nginx:1.25-alpine
EOF
```

### Step 4: Verify from your Host Terminal
```bash
kubectl get pods -n default
# Expected output: my-static-web-<node-name>   1/1   Running
```
