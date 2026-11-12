#!/bin/bash
set -euo pipefail

cat << 'EOF' > /usr/local/bin/backup-tool
#!/bin/bash
echo "Executing system maintenance backup..."
EOF

chown root:root /usr/local/bin/backup-tool
# Insecure permissions: SUID (4) + world writable (777)
chmod 4777 /usr/local/bin/backup-tool
