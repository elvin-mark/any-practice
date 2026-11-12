#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

cat <<'EOF' > app.sh
#!/bin/sh
echo "Running as UID $(id -u)"
echo "Writing timestamp..." > /data/timestamp.txt
cat /data/timestamp.txt
EOF
chmod +x app.sh

cat <<'EOF' > Dockerfile
FROM alpine:3.20
WORKDIR /app
RUN mkdir -p /app /data
COPY app.sh /app/app.sh
# Flawed: runs as root by default
CMD ["/app/app.sh"]
EOF

echo "[PASS] Environment prepared for img-03-security-nonroot in $DIR"
