#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

cat <<'EOF' > package.json
{
  "name": "sample-node-service",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.19.2"
  }
}
EOF

cat <<'EOF' > server.js
console.log("App server running v1.0");
EOF

cat <<'EOF' > Dockerfile
FROM alpine:3.20
WORKDIR /app
# Flawed ordering: copies everything first, busting cache every time server.js changes
COPY . .
RUN echo "Installing deps..." && touch /app/deps_installed
CMD ["cat", "/app/deps_installed"]
EOF

echo "[PASS] Environment prepared for img-02-layer-caching in $DIR"
