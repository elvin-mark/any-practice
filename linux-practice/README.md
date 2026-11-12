# Linux Practice Platform (`lp`)

`lp` (`linux-practice`) is an open-source, local CLI and interactive TUI platform designed to help engineers, students, and sysadmins master Linux systems administration hands-on.

Instead of passive reading, `lp` provisions ephemeral, isolated **Ubuntu 24.04** containers running full **systemd** init, injects real-world failures, configuration errors, and administrative challenges, and tests your solutions using automated verification harnesses.

---

## 🌟 Key Features

- **🐧 Real Systemd Environments**: Every lab runs full `systemd` init (not just mock processes or stripped containers), allowing authentic practice with `systemctl`, `journalctl`, systemd timers, and ufw firewalls.
- **⚡ Dual Interaction Mode**:
  - **Interactive TUI**: Rich terminal interface (built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lipgloss](https://github.com/charmbracelet/lipgloss)) with keyboard navigation, progress statistics, and side-by-side challenge descriptions.
  - **Scriptable CLI**: Direct commands (`lp start`, `lp verify`, `lp hint`, etc.) for seamless integration with your favorite terminal, tmux, or automated testing.
- **📦 Self-Contained & Offline-Ready**: Ships with embedded base systemd images and curated core scenarios compiled directly into a single Go binary (`go:embed`). No external network connection required once built.
- **📊 Local Progress Tracking**: Automatically tracks attempts, completion status, and hints viewed in a human-readable state file (`~/.config/linux-practice/state.yaml`).
- **🛠️ Scenario Authoring Tools**: Easily create and test community scenarios with `lp create` and `lp validate`.

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- Docker or Podman running locally (supports Docker Engine, Docker Desktop for Linux, and rootless Docker/Podman)

### 1. Build the Binary
```bash
git clone https://github.com/elvin-mark/linux-practice.git
cd linux-practice
make build
```

The compiled binary will be placed at `./bin/lp`. You can add it to your `PATH` or install it globally with `make install` or `go install ./cmd/lp`.

### 2. Launch the Interactive TUI
Run `lp` without arguments:
```bash
./bin/lp
```

### 3. Or Use the CLI Directly
```bash
# List all available scenarios
./bin/lp list

# Start a scenario and attach an interactive shell
./bin/lp start perm-suid-cleanup

# If you exited the shell, verify your fix
./bin/lp verify perm-suid-cleanup

# Need a pointer?
./bin/lp hint perm-suid-cleanup

# Re-attach to the active container shell
./bin/lp attach

# View your overall progress
./bin/lp progress

# Tear down the scenario container
./bin/lp stop perm-suid-cleanup
```

---

## ⌨️ TUI Keybindings

| Key | Action |
| :--- | :--- |
| `↑` / `k`, `↓` / `j` | Navigate scenario list |
| `Enter` / `s` | Start selected scenario and open shell |
| `a` | Re-attach to active container shell |
| `v` | Verify current scenario solution |
| `h` | Reveal next hint |
| `r` | Reset scenario to starting state |
| `q` / `Ctrl+C` | Quit TUI |

---

## 📚 Included Starter Scenarios

| Scenario ID | Track | Difficulty | Description |
| :--- | :--- | :--- | :--- |
| **`perm-suid-cleanup`** | Permissions | Beginner | Audit `/usr/local/bin` to find and revoke dangerous SUID and world-writable permissions. |
| **`user-sudo-syntax`** | Users | Beginner | Diagnose and fix a syntax error in `/etc/sudoers.d/devops` using `visudo`. |
| **`systemd-failed-service`**| Systemd | Intermediate | Inspect `journalctl` to diagnose a failing service unit with directory permission issues. |
| **`systemd-custom-timer`**  | Systemd | Intermediate | Replace a legacy cron script by authoring a native systemd timer unit scheduled every 15 minutes. |
| **`net-ufw-lockdown`**      | Networking | Intermediate | Configure UFW with default deny incoming, default allow outgoing, and allow SSH + HTTP. |
| **`troubleshoot-port-conflict`** | Troubleshooting | Advanced | Diagnose an incident where port 80 is blocked by a rogue process, terminate it, and restore the web service. |

---

## 🛠️ Authoring Custom Scenarios

You can create your own scenarios for practice or contribute them back:

### 1. Scaffold a Scenario Directory
```bash
./bin/lp create my-custom-scenario
```
This generates:
```
my-custom-scenario/
├── scenario.yaml   # Metadata, track, tags, container settings, hints
├── setup.sh        # Bash script executed inside container to inject the problem
└── verify.sh       # Bash test script returning exit code 0 on success
```

### 2. Validate the Scenario
```bash
./bin/lp validate ./my-custom-scenario
```

### 3. Add to User Scenarios
Drop any valid scenario folder into `~/.config/linux-practice/scenarios/`. `lp` will automatically discover and load it alongside the embedded core curriculum!

---

## 🧪 Testing

Run the full automated test suite:
```bash
make test
```

---

## 📄 License
MIT License.
