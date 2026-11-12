# 🐳 Docker Practice Platform (`dockp`)

> **An interactive, terminal-native training platform for mastering Docker containerization, production architecture, and security hardening.**

`dockp` (`docker-practice`) is part of the `any-practice` suite. It provides hands-on, authentic lab scenarios designed to take you from core container concepts to production-grade Docker and DCA (Docker Certified Associate) expertise.

Instead of passive tutorials, `dockp` injects authentic container challenges, broken architectures, and security misconfigurations, validating your solutions automatically using real Docker engine assertions.

---

## 🌟 Key Features

- **🐳 Authentic Docker Scenarios**: Practice multi-stage image builds, cgroup resource limits, custom bridge topologies, volume backup pipelines, Docker Compose dependency ordering, and Linux capability dropping.
- **⚡ Dual Interaction Modes**:
  - **Interactive TUI**: Split-pane [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss) dashboard with progress stats, domain filters, and real-time verification.
  - **Scriptable CLI**: Direct commands (`dockp list`, `dockp start`, `dockp verify`, `dockp hint`, `dockp solution`, `dockp reset`) for fast terminal workflows and scripting.
- **🎯 7 Core Domains & 17 Hands-on Labs**:
  1. **Dockerfile & Image Optimization**: Multi-stage builds, layer caching order, unprivileged non-root users.
  2. **Container Lifecycle & Resources**: CPU/memory limits, healthcheck probes, PID 1 zombie reaping (`--init`).
  3. **Storage, Volumes & Persistence**: Named volumes, read-only root filesystems with tmpfs, backup/restore pipelines.
  4. **Networking & Service Discovery**: User-defined bridge DNS discovery, loopback IP bindings, DMZ vs backend network isolation.
  5. **Docker Compose Orchestration**: Multi-container stacks, health-based startup dependencies (`service_healthy`).
  6. **Container Security & Hardening**: Linux capability dropping (`--cap-drop=ALL`), disabling setuid escalation (`no-new-privileges`).
  7. **Troubleshooting & Diagnostics**: OOMKilled exit code 137 analysis, broken entrypoints with missing shebangs/permissions.
- **📊 Persistent Progress Tracking**: Automatically tracks completion status, points earned, hints revealed, and attempts in `~/.dockp/progress.json`.
- **🧹 Clean Lab Environment**: Built-in `dockp clean` cleans up practice containers, networks, and volumes tagged with `practice=docker`.

---

## 🚀 Quick Start

### Prerequisites
- **Go**: 1.22+ (configured with `go.work`)
- **Docker Engine**: Docker 20.10+ / Docker 24+ / Docker 29+
- **Docker Compose**: v2 plugin (`docker compose`)

### 1. Build
```bash
# From docker-practice directory
make build

# Or from root any-practice directory
make build-dockp
```
The compiled binary will be placed at `./bin/dockp`.

### 2. Check Prerequisites
```bash
./bin/dockp check
```

### 3. Launch Interactive TUI
```bash
./bin/dockp
# or
./bin/dockp tui
```

### 4. CLI Usage
```bash
# List all practice scenarios
./bin/dockp list

# Start a specific scenario lab
./bin/dockp start img-01-multistage-build

# Read instructions and complete the task in the scenario directory
cd scenarios/img-01-multistage-build

# Verify your fix against the Docker engine
dockp verify

# Reveal progressive hints if stuck
dockp hint

# View the official reference solution
dockp solution img-01-multistage-build

# View overall progress and domain mastery
dockp progress

# Clean up practice containers and networks
dockp clean
```

---

## ⌨️ TUI Keybindings

| Key | Action |
| :--- | :--- |
| `↑` / `k`, `↓` / `j` | Navigate scenario list |
| `Enter` / `p` | Start / practice selected scenario |
| `v` | Verify current scenario against Docker engine |
| `h` | Reveal next hint |
| `s` | Toggle official reference solution guide |
| `r` | Reset scenario to initial starting state |
| `q` / `Ctrl+C` | Quit TUI |

---

## 📂 Scenario Curriculum

| ID | Title | Domain | Difficulty | Points |
| :--- | :--- | :--- | :--- | :--- |
| `img-01-multistage-build` | Multi-Stage Build for Go Microservice | Images | Medium | 7 |
| `img-02-layer-caching` | Optimize Dockerfile Layer Caching | Images | Easy | 5 |
| `img-03-security-nonroot` | Harden Dockerfile Non-Root User | Images | Medium | 6 |
| `run-01-resource-limits` | Enforce CPU, Memory & Restart Policies | Lifecycle | Easy | 5 |
| `run-02-healthchecks` | Implement Container Healthchecks | Lifecycle | Medium | 6 |
| `run-03-init-zombies` | Fix Zombie Reaping with Docker Init | Lifecycle | Medium | 7 |
| `vol-01-named-volume` | Persistent Storage with Named Volumes | Storage | Easy | 5 |
| `vol-02-readonly-tmpfs` | Read-Only Rootfs with Writable Tmpfs | Storage | Medium | 6 |
| `vol-03-backup-restore` | Volume Backup & Restoration Pipeline | Storage | Hard | 8 |
| `net-01-custom-bridge` | User-Defined Bridge Network & DNS | Networking | Easy | 5 |
| `net-02-port-binding` | Secure Host IP Port Binding | Networking | Medium | 6 |
| `net-03-isolated-backend` | Multi-Tier Network Isolation (DMZ) | Networking | Hard | 8 |
| `comp-01-web-db-stack` | Compose Stack with Health Dependencies | Compose | Medium | 7 |
| `sec-01-cap-drop` | Least Privilege: Capability Dropping | Security | Medium | 7 |
| `sec-02-no-new-privs` | Prevent SUID Escalation | Security | Easy | 5 |
| `trouble-01-oom-killer` | Diagnose & Fix OOMKilled Container | Trouble | Medium | 7 |
| `trouble-02-broken-entrypoint` | Fix Exec Format & Missing Shebang | Trouble | Medium | 6 |

---

## 📜 License

MIT License. Designed and maintained by [Elvin Mark](https://github.com/elvin-mark).
