#!/usr/bin/env bash
set -e

CP_CONTAINER=$(docker ps --filter "name=control-plane" --format "{{.Names}}" | head -n 1)
if [ -z "$CP_CONTAINER" ]; then
  echo "Control-plane container not found"
  exit 1
fi

# Ensure etcdctl wrapper is available
docker exec "$CP_CONTAINER" sh -c 'cat << "EOF" > /usr/local/bin/etcdctl
#!/bin/sh
ETCD_ID=$(crictl ps --name etcd -q | head -n 1)
exec crictl exec $ETCD_ID etcdctl "$@"
EOF
chmod +x /usr/local/bin/etcdctl'

# Clean previous backup
docker exec "$CP_CONTAINER" rm -rf /var/lib/etcd-backup 2>/dev/null || true
docker exec "$CP_CONTAINER" mkdir -p /var/lib/etcd-backup

echo "Prepared control-plane node with etcdctl utility and clean /var/lib/etcd-backup directory."
