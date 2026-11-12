# Solution: Troubleshoot Control Plane Static Pod Failure

### Step 1: Connect to Control Plane Node
```bash
docker exec -it <control-plane-container> bash
```

### Step 2: Inspect Static Pod Manifests
```bash
cd /etc/kubernetes/manifests
cat infra-agent.yaml
```
Notice the invalid command: `command: ["invalid-command-does-not-exist"]`.

### Step 3: Edit the Manifest
Edit `/etc/kubernetes/manifests/infra-agent.yaml` to provide a valid command:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: infra-agent
  namespace: kube-system
spec:
  restartPolicy: Always
  containers:
  - name: agent
    image: busybox:1.36
    command: ["sh", "-c", "echo 'Infra agent running'; sleep 3600"]
```

### Step 4: Verify the Pod Runs
```bash
kubectl get pods -n kube-system | grep infra-agent
# Status should transition to Running
```
