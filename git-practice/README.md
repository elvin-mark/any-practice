# git-practice (gitp)

> Interactive Git mastery, history surgery, conflict resolution, and disaster recovery drills.

`gitp` provides safe, automated hands-on practice labs for Git workflows:
- **Interactive Rebase**: Squashing WIP commits, editing history, fixing commit messages.
- **Merge Conflicts**: Resolving real 3-way branch conflicts cleanly.
- **Reflog Recovery**: Rescuing orphaned commits and deleted branches using `git reflog`.
- **Git Bisect**: Binary-search bug hunting across commit histories with automated test scripts.
- **Branching & Detached HEAD**: Cherry-picking patches and rescuing work from detached HEAD states.

Each drill creates an isolated Git repository in `/tmp/gitp-labs/<scenario-id>` where you perform real Git commands, followed by automated verification.

---

## Quickstart

```bash
# Build binary
make build

# Launch interactive split-pane TUI
./bin/gitp

# List all available drills
./bin/gitp list

# Start a specific lab
./bin/gitp start rebase-01-squash

# Switch into the lab directory and work:
cd /tmp/gitp-labs/rebase-01-squash
git rebase -i HEAD~3

# Verify your solution
gitp verify rebase-01-squash

# View progress & skill breakdown
./bin/gitp progress
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
id: my-scenario
title: Scenario Name
category: Rebase
difficulty: Easy
description: Problem statement...
goal: Expected end state...
hints:
  - First hint
setup: |
  # Shell commands to construct initial git repo
  git commit -m "initial"
verify: |
  # Shell commands returning 0 on pass, non-zero on failure
  git log -1 --pretty=%s | grep -q "expected commit"
```
