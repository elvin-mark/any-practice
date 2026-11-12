#!/bin/bash
set -euo pipefail

# Ensure UFW is installed and in a clean disabled state
ufw --force reset
ufw disable || true
