# 🥋 Any-Practice

> **A unified, terminal-native practice suite for mastering Software Engineering, Linux Systems Administration, Kubernetes (CKA), Docker Containerization, and Capture The Flag (CTF) Cybersecurity.**

`any-practice` brings together five specialized, interactive CLI/TUI tools into a single cohesive monorepo. Built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss), all tools share a unified visual design language, dual CLI/TUI interfaces, and automated validation pipelines.

---

## 🛠️ The Practice Suite

| Tool | Focus | Interactive Mode | CLI Mode | Target Environments |
| :--- | :--- | :--- | :--- | :--- |
| **`algo`** | Computer Science & Algorithms | `algo` / `algo watch` | `algo list`, `algo verify <id>` | Pure Go, AST validation, benchmarks |
| **`lp`** | Linux Systems Administration | `lp` | `lp list`, `lp start <id>`, `lp verify` | Ubuntu 24.04 systemd containers |
| **`cka`** | Certified Kubernetes Administrator | `cka` / `cka tui` | `cka list`, `cka practice <id>`, `cka exam` | Local Kind multi-node clusters + Calico |
| **`dockp`** | Docker Containerization & Security | `dockp` / `dockp tui` | `dockp list`, `dockp start <id>`, `dockp verify` | Production Docker Engine & Compose |
| **`ctf`** | Cybersecurity & Capture The Flag | `ctf` / `ctf tui` | `ctf list`, `ctf submit <id> <flag>` | Security puzzles, crypto, web, pwn, rev, forensics |

---

## 🎨 Unified Design System

Every tool in the suite adheres to the **Any-Practice Terminal Design System**:

- **Color Palette**:
  - 🟣 **Primary Brand**: `#7D56F4` (Electric Violet) & `#5A32A3` (Deep Purple)
  - 🔷 **Accent / Highlight**: `#00D7D7` (Cyan)
  - 🟢 **Success**: `#04B575` (Emerald Green)
  - 🟡 **Warning / Points**: `#FFAE00` (Amber Gold)
  - 🔴 **Error / Failures**: `#FF4040` (Coral Red)
  - ⬛ **Cards & Borders**: `#3C3854` (Muted Purple-Slate) with rounded box borders
- **Consistent Visual Identity**:
  - **Tool Banners**: Unified header pill with tool branding, subtitle, and subtle horizontal rule.
  - **Unicode Tabular Output**: Aligned `list` commands rendering Unicode box tables (`┌─┬─┐`) with status checks and difficulty badges.
  - **Interactive Split-Pane TUIs**: Dual-panel Bubbletea interface (category list on left, scenario/problem details card on right).
  - **Standardized Messaging**: Color-coded prefixes (`✔ PASS:`, `✖ FAIL:`, `ℹ INFO:`, `⚠ WARN:`, `💡 HINT:`).
  - **Keybinding Footers**: Key badges in bold violet with subtle action descriptions (`↑/↓ navigate • enter start • v verify • q quit`).

---

## 🚀 Quick Start

### Prerequisites
- **Go**: 1.22+ (configured with `go.work`)
- **Docker** or **Podman**: For `lp` (Linux containers) and `cka` (Kind Kubernetes nodes)
- **kubectl** & **kind**: For `cka` practice labs

### 1. Build All Tools
```bash
make build
```
Binaries will be placed in `./bin/`:
- `bin/algo`
- `bin/lp`
- `bin/cka`
- `bin/dockp`
- `bin/ctf`

Run tests across all five modules:
```bash
make test
```

---

## 📖 Tool Usage

### ⚡ 1. Algorithm Practice (`algo`)
Hands-on Go algorithm drills with live watcher, AST validation (anti-cheating standard library checks), and execution benchmarks.

```bash
# Launch interactive TUI watcher (watches exercises/ directory on save)
./bin/algo

# List curriculum with completion status
./bin/algo list

# Initialize / unpack exercises into current directory
./bin/algo init

# Verify a specific exercise from CLI
./bin/algo verify sorting_01_bubble_sort

# View hints or solution
./bin/algo hint sorting_01_bubble_sort
./bin/algo solution sorting_01_bubble_sort
```

---

### 🐧 2. Linux Practice Platform (`lp`)
Provisions authentic systemd containers, injects real-world failures, and tests your fixes automatically.

```bash
# Launch interactive TUI scenario browser & runner
./bin/lp

# List available Linux administration scenarios
./bin/lp list

# Start a scenario and attach to interactive shell (as student with sudo)
./bin/lp start net-ufw-lockdown

# Verify your fix
./bin/lp verify

# View hints or check track progress
./bin/lp hint
./bin/lp progress
```

---

### ☸ 3. CKA Kubernetes Practice (`cka`)
Simulate the 2-hour Certified Kubernetes Administrator exam or drill individual domain tasks using local Kind clusters.

```bash
# Launch interactive CKA scenario browser & lab dashboard
./bin/cka

# Initialize local Kind clusters with Calico CNI
./bin/cka cluster init

# List all 20 CKA scenarios across the 5 CNCF exam domains
./bin/cka list

# Practice an individual scenario
./bin/cka practice arch-01-rbac

# Verify scenario solution against the cluster
./bin/cka verify

# Start a full 2-hour timed mock exam simulation
./bin/cka exam start
./bin/cka exam status
./bin/cka exam next
./bin/cka exam finish
```

---

### 🐳 4. Docker Practice Platform (`dockp`)
Hands-on Docker containerization, multi-stage image optimization, networking, persistent volume pipelines, Compose stacks, and security hardening.

```bash
# Launch interactive Docker TUI browser
./bin/dockp

# Check local Docker engine prerequisites
./bin/dockp check

# List all 17 scenarios across 7 core domains
./bin/dockp list

# Start a practice scenario lab
./bin/dockp start img-01-multistage-build

# Verify solution against Docker engine
./bin/dockp verify

# View hints or official solution guide
./bin/dockp hint
./bin/dockp solution img-01-multistage-build

# View overall progress and domain mastery
./bin/dockp progress

# Clean up practice lab containers and networks
./bin/dockp clean
```

---

### 🚩 5. Capture The Flag Practice (`ctf`)
Interactive CTF training platform covering cryptography, web security, binary reverse engineering, stack overflows, forensics, and decoding drills.

```bash
# Launch interactive split-pane TUI dashboard
./bin/ctf

# List all available challenges and points
./bin/ctf list

# Display challenge briefing and instructions
./bin/ctf show crypto-01-caesar

# Reveal progressive hints
./bin/ctf hint crypto-01-caesar

# Submit a flag
./bin/ctf submit crypto-01-caesar "FLAG{caesar_shift_is_classic}"

# View category breakdown and total score
./bin/ctf progress
```

---

## 📂 Repository Structure

```
any-practice/
├── Makefile                    # Monorepo build, test, and cleanup targets
├── go.work                     # Go workspace linking all four modules
├── README.md                   # Suite overview and documentation
├── bin/                        # Compiled suite binaries (algo, lp, cka, ctf)
├── algo-practice/              # Go algorithm practice workbook
│   ├── cmd/algo/               # algo CLI entrypoint
│   ├── exercises/              # Practice exercises across 5 tracks
│   ├── internal/ui/            # Unified design system & CLI table renderer
│   └── internal/tui/           # Bubble Tea live watcher TUI
├── linux-practice/             # Linux systemd hands-on labs
│   ├── cmd/lp/                 # lp CLI entrypoint
│   ├── embedded/               # Embedded scenarios & systemd Dockerfile
│   ├── internal/ui/            # Unified design system & CLI table renderer
│   ├── internal/cli/           # Cobra CLI commands
│   └── internal/tui/           # Bubble Tea split-pane TUI
├── k8s-practice/               # CKA mock exam simulator & drill labs
│   ├── cmd/cka/                # cka CLI entrypoint
│   ├── scenarios/              # 20 CKA scenario YAMLs, prompts, and verification scripts
│   ├── internal/ui/            # Unified design system & CLI table renderer
│   ├── internal/tui/           # Bubble Tea CKA scenario browser TUI
│   ├── internal/cluster/       # Kind cluster management
│   └── internal/exam/          # Exam timer, state, and score report generator
├── docker-practice/            # Docker containerization & security drills
│   ├── cmd/dockp/              # dockp CLI entrypoint
│   ├── scenarios/              # 17 Docker scenario YAMLs, problems, setup/verify scripts
│   ├── internal/ui/            # Unified design system & CLI table renderer
│   ├── internal/tui/           # Bubble Tea split-pane TUI
│   ├── internal/scenario/      # Scenario model, loader, and automated runner
│   └── internal/progress/      # Progress tracker and scoring engine
└── ctf-practice/               # Capture The Flag security training
    ├── cmd/ctf/                # ctf CLI entrypoint
    ├── challenges/             # Curated challenges across 6 categories
    ├── internal/ui/            # Unified design system & CLI table renderer
    ├── internal/tui/           # Bubble Tea interactive dashboard
    ├── internal/challenge/     # Challenge model, registry, and flag validation
    └── internal/progress/      # Uniform progress and score tracker
```

---

## 📜 License

MIT License. Designed and maintained by [Elvin Mark](https://github.com/elvin-mark).
