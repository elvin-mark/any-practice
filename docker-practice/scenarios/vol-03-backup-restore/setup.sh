#!/usr/bin/env bash
set -e

rm -f /tmp/prod-backup.tar.gz || true
docker volume rm -f prod-backup-restored >/dev/null 2>&1 || true
docker volume rm -f prod-backup-source >/dev/null 2>&1 || true

# Provision source volume with data
docker volume create --label practice=docker prod-backup-source >/dev/null

docker run --rm -v prod-backup-source:/data alpine:3.20 sh -c '
cat <<EOF > /data/config.json
{"cluster": "production-us-east", "replicas": 3, "tls": true}
EOF
echo "SQL_DUMP_2026_MASTER_DATA_V99" > /data/data.db
sha256sum /data/config.json /data/data.db > /data/checksum.sha256
'

echo "[PASS] Environment prepared for vol-03-backup-restore with source volume prod-backup-source"
