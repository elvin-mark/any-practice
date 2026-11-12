# Solution: Take an etcd Snapshot Backup on Control Plane

### Step 1: Connect to the Control-Plane Node
```bash
docker exec -it <control-plane-container> bash
```

### Step 2: Inspect etcd Pod Manifest for Certificates
```bash
grep -E 'cert-file|key-file|trusted-ca-file|listen-client-urls' /etc/kubernetes/manifests/etcd.yaml
```
Output reveals:
- `--cacert=/etc/kubernetes/pki/etcd/ca.crt`
- `--cert=/etc/kubernetes/pki/etcd/server.crt`
- `--key=/etc/kubernetes/pki/etcd/server.key`
- `--endpoints=https://127.0.0.1:2379`

### Step 3: Run the etcdctl Snapshot Save Command
```bash
mkdir -p /var/lib/etcd-backup

ETCDCTL_API=3 etcdctl \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key \
  snapshot save /var/lib/etcd-backup/snapshot.db
```

### Step 4: Verify the Snapshot
```bash
ETCDCTL_API=3 etcdctl --write-out=table snapshot status /var/lib/etcd-backup/snapshot.db
```
