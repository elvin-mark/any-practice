#!/bin/bash
set -euo pipefail

# Ensure devops group exists
getent group devops >/dev/null || groupadd devops

# Intentionally misspelled NOPASSWD as NOPASWD
cat << 'EOF' > /etc/sudoers.d/devops
%devops ALL=(ALL:ALL) NOPASWD: ALL
EOF

chmod 0440 /etc/sudoers.d/devops
