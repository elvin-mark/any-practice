#!/bin/bash
set -euo pipefail

# Ensure target working directory is missing
rm -rf /var/data/worker

cat << 'EOF' > /usr/local/bin/data-worker.sh
#!/bin/bash
echo "Worker started at $(date)" >> worker.log
while true; do
  sleep 5
done
EOF

chmod +x /usr/local/bin/data-worker.sh

cat << 'EOF' > /etc/systemd/system/data-worker.service
[Unit]
Description=Data Processing Worker Service
After=network.target

[Service]
Type=simple
User=student
Group=student
WorkingDirectory=/var/data/worker
ExecStart=/usr/local/bin/data-worker.sh
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
# Start it so it enters failed state
systemctl start data-worker.service || true
