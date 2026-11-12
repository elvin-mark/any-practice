#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

docker rm -f fixed-worker >/dev/null 2>&1 || true

# Write broken entrypoint (missing shebang, not executable)
printf "echo 'ENTRYPOINT_EXECUTION_SUCCESSFUL'\nsleep 3600\n" > entrypoint.sh
chmod -x entrypoint.sh

cat <<'EOF' > Dockerfile
FROM alpine:3.20
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
EOF

echo "[PASS] Environment prepared for trouble-02-broken-entrypoint in $DIR"
