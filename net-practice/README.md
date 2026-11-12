# net-practice (netp)

> Hands-on Computer Networking, Protocols, Packet Traces, Routing & Port Conflict Labs.

`netp` provides safe, automated hands-on practice labs for computer networking:
- **DNS Troubleshooting**: Fixing resolver loops, missing upstream nameservers, `/etc/resolv.conf`, and `hosts` overrides.
- **CIDR Subnetting**: Network prefix boundaries, broadcast addresses, and usable host address calculations.
- **TLS Certificate Bundling**: Inspecting certificate chains with OpenSSL, resolving missing intermediate CA errors.
- **Packet & Traffic Analysis**: Parsing HTTP and TCP packet logs to extract leaked authorization tokens.
- **Routing Table Gateways**: Configuring default gateways and static routes via standard Linux route semantics.
- **Sockets & Port Conflicts**: Diagnosing `address already in use` socket conflicts and terminating rogue listeners.

---

## Quickstart

```bash
# Build binary
make build

# Launch interactive split-pane TUI
./bin/netp

# List all available networking drills
./bin/netp list

# Start a specific lab (sets up /tmp/netp-labs/<id>)
./bin/netp start dns-01-resolve-failure

# Switch into the lab directory and work:
cd /tmp/netp-labs/dns-01-resolve-failure

# Verify your solution
netp verify dns-01-resolve-failure

# View progress & skill breakdown
./bin/netp progress
```

---

## Interactive TUI Controls

| Key | Action |
| --- | --- |
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `Enter` / `s` | Initialize & start lab |
| `v` | Verify lab solution |
| `r` | Reset workspace to initial state |
| `h` | Reveal next progressive hint |
| `q` / `Ctrl+C` | Quit |

---

## Adding New Scenarios

Add new YAML scenario definitions into `scenarios/`:

```yaml
id: my-net-drill
title: "Drill Name"
category: DNS
difficulty: Easy
description: Problem statement...
goal: Expected end state...
hints:
  - First hint
setup: |
  # Shell commands to construct initial lab environment
verify: |
  # Shell commands returning 0 on pass, non-zero on failure
  grep -q "expected" config
```
