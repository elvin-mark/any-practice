#!/bin/bash
set -euo pipefail

cat << 'EOF' > /usr/local/bin/backup-report.sh
#!/bin/bash
echo "Backup report generated: $(date)" >> /var/log/backup-report.log
EOF
chmod +x /usr/local/bin/backup-report.sh

cat << 'EOF' > /etc/systemd/system/backup-report.service
[Unit]
Description=Generate Backup Report

[Service]
Type=oneshot
ExecStart=/usr/local/bin/backup-report.sh
EOF

# Ensure no stale timer exists
rm -f /etc/systemd/system/backup-report.timer
systemctl daemon-reload
