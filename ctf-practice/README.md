# ctf-practice (ctf)

> Capture The Flag (CTF) security training platform and vulnerability challenges.

`ctf` provides hands-on practice across core cybersecurity and CTF categories:
- **Crypto**: Substitution ciphers, XOR schemes, modular arithmetic, RSA basics.
- **Web**: SQL injection, authentication bypass, session manipulation, XSS.
- **Forensics**: Metadata extraction, file carving, steganography, memory dumps.
- **Reverse**: String discovery, disassembled logic, binary reversing.
- **Pwn**: Stack buffer overflows, memory corruption concepts, return address hijacking.
- **Misc**: Multi-layer encodings, esoteric formats, puzzle solving.

---

## Quickstart

```bash
# Build binary
make build

# Launch interactive split-pane TUI
./bin/ctf

# List all challenges with category, points, and completion status
./bin/ctf list

# View briefing for a specific challenge
./bin/ctf show crypto-01-caesar

# Reveal progressive hints
./bin/ctf hint crypto-01-caesar

# Submit a flag
./bin/ctf submit crypto-01-caesar "FLAG{caesar_shift_is_classic}"

# View progress breakdown and scoreboard
./bin/ctf progress
```

---

## Interactive TUI Controls

| Key | Action |
| --- | --- |
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `s` / `Enter` | Submit flag prompt |
| `h` | Reveal next hint |
| `Esc` | Cancel submission prompt |
| `q` / `Ctrl+C` | Quit |

---

## Adding New Challenges

Place new YAML or JSON files into `challenges/`:

```yaml
id: web-02-cookie-tamper
title: Cookie Tampering
category: Web
difficulty: Medium
points: 100
description: |
  Inspect the base64-encoded cookie and elevate privileges to 'admin'.
hints:
  - "Decode the cookie: echo $COOKIE | base64 -d"
  - "Change role=user to role=admin and re-encode"
flag: FLAG{tamper_proof_signatures_matter}
```
