#!/bin/bash
set -euo pipefail

# Ensure curl is installed
apt-get install -y curl

# Spawn a rogue python HTTP server listening on port 80 in the background
nohup python3 -m http.server 80 > /var/log/rogue-server.log 2>&1 &

cat << 'EOF' > /var/www/index.html
<html><body><h1>Production Web Service Active</h1></body></html>
EOF

cat << 'EOF' > /etc/systemd/system/web-app.service
[Unit]
Description=Production Web Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/var/www
ExecStart=/usr/bin/python3 -m http.server 80
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
# Attempt to start so systemd records failure
systemctl start web-app.service || true
