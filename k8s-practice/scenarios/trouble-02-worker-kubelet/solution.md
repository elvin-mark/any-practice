# Solution: Troubleshoot Worker Node in NotReady Status

### Step 1: Identify NotReady Node
```bash
kubectl get nodes
# Identify the node name, e.g. dev-cluster-worker or cka-cluster-1-worker
```

### Step 2: Describe the Node for Clues
```bash
kubectl describe node <worker-node-name>
# Under Conditions, check Kubeletstopped or KubeletNotReady
```

### Step 3: Access the Node Container
```bash
docker exec -it <worker-node-name> bash
```

### Step 4: Check Kubelet Service Status
```bash
systemctl status kubelet
journalctl -u kubelet -n 30
```

### Step 5: Start / Restart Kubelet Service
```bash
systemctl daemon-reload
systemctl restart kubelet
systemctl enable kubelet
```

### Step 6: Verify from Host
```bash
kubectl get nodes
# Node should transition to Ready within 10-30 seconds.
```
