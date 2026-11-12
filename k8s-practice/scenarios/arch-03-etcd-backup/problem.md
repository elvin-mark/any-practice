# Task: Take an etcd Snapshot Backup on Control Plane

You are required to perform an etcd backup of the Kubernetes cluster.

1. Access the control plane node (e.g. `docker exec -it <control-plane-node> bash`).
2. Identify the etcd cluster endpoint, CA certificate, client certificate, and client key from `/etc/kubernetes/manifests/etcd.yaml`.
3. Take a snapshot of etcd using `etcdctl` with API version 3.
4. Save the snapshot file to `/var/lib/etcd-backup/snapshot.db` on the control plane node.
5. Verify the snapshot integrity using `etcdctl snapshot status`.
