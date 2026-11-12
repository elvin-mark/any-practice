# ⚡ `algo` — Interactive Go Algorithm Practice Platform

A terminal-native, interactive algorithm workbook and test runner built in Go. Inspired by tools like *Rustlings* and *Go-tour*, `algo` helps you master essential computer science algorithms by writing real Go code in your favorite editor while a live TUI continuously watches, tests, benchmarks, and lints your solutions.

---

## 🚀 Key Features

- **Continuous TUI Watcher:** Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss). Automatically detects file saves and displays live feedback in real time.
- **Three-Tier Validation Pipeline:**
  1. **AST Linter (`go/parser`, `go/ast`):** Disallows standard library shortcuts (e.g. banning `sort`, `slices.Sort`, `container/heap`) to ensure you implement the underlying algorithm logic from scratch.
  2. **Comprehensive Unit Tests (`go test`):** Rigorous edge cases (empty slices, duplicates, single elements, negative numbers, large datasets) with timeouts to catch infinite loops.
  3. **Complexity & Allocation Benchmarks:** Measures runtime (`ns/op`), heap memory (`B/op`), and allocations (`allocs/op`) to enforce optimal Big-O complexity.
- **3-Tier Progressive Hint System:**
  - **Tier 1:** Conceptual approach & intuition.
  - **Tier 2:** Invariants & tricky edge cases.
  - **Tier 3:** Algorithmic pseudocode walkthrough.
- **Opt-in Reference Solutions:** View canonical, idiomatic Go implementations directly in the TUI (with a confirmation safeguard for unsolved exercises).
- **Curriculum Browser:** Seamlessly jump across tracks or navigate sequentially.
- **Progress Persistence:** Automatic tracking in `.algo/progress.json` with best benchmark timings and attempt counters.
- **Dual Distribution:** Run directly from the cloned repository or install globally via `go install` with `algo init` unpacking embedded exercises.

---

## 📦 Getting Started

### 1. Run Directly (Cloned Repository)

Clone the repository and run:

```bash
# Start the interactive watcher
go run ./cmd/algo watch

# Or compile the binary once
go build -o algo ./cmd/algo
./algo watch
```

### 2. Install Globally

```bash
go install github.com/elvin-mark/algo-practice/cmd/algo@latest

# In any empty directory of your choice:
algo init

# Launch the practice platform
algo watch
```

---

## ⌨️ TUI Keybindings

While `algo watch` is running:

| Key | Action |
|:---:|:---|
| `[h]` | Cycle / reveal next progressive hint (1 -> 2 -> 3) |
| `[s]` | View reference solution (prompts confirmation if unsolved) |
| `[n]` | Move to next exercise |
| `[p]` | Move to previous exercise |
| `[m]` | Open curriculum browser / track picker |
| `[r]` | Force re-run verification |
| `[q]` / `Ctrl+C` | Quit |

---

## 📚 Seed Curriculum Tracks (15 Core Algorithms)

### 1. 📂 Sorting & Searching
1. **Bubble Sort** (`sorting_01_bubble_sort`) — Easy — *O(N^2) Time, O(1) Space*
2. **Binary Search** (`sorting_02_binary_search`) — Easy — *O(log N) Time, O(1) Space*
3. **Merge Sort** (`sorting_03_merge_sort`) — Medium — *O(N log N) Time, O(N) Space*
4. **Quick Sort** (`sorting_04_quick_sort`) — Medium — *O(N log N) Time, O(log N) Stack*

### 2. 📂 Linked Lists
1. **Reverse Linked List** (`linked_list_01_reverse`) — Easy — *O(N) Time, O(1) Space*
2. **Detect Cycle (Floyd's Tortoise & Hare)** (`linked_list_02_detect_cycle`) — Easy — *O(N) Time, O(1) Space*
3. **Merge Two Sorted Lists** (`linked_list_03_merge_two_sorted`) — Easy — *O(N + M) Time, O(1) Space*

### 3. 📂 Stacks & Queues
1. **Valid Parentheses** (`stack_01_valid_parentheses`) — Easy — *O(N) Time, O(N) Space*
2. **Min Stack** (`stack_02_min_stack`) — Medium — *O(1) All operations, O(N) Space*
3. **Queue Using Two Stacks** (`stack_03_queue_using_stacks`) — Easy — *Amortized O(1), O(N) Space*

### 4. 📂 Trees
1. **In-Order Traversal** (`tree_01_inorder_traversal`) — Easy — *O(N) Time, O(N) Space*
2. **Invert Binary Tree** (`tree_02_invert_binary_tree`) — Easy — *O(N) Time, O(H) Space*
3. **Maximum Depth of Binary Tree** (`tree_03_max_depth`) — Easy — *O(N) Time, O(H) Space*

### 5. 📂 Two Pointers & Sliding Window
1. **Two Sum (Sorted Array)** (`twopointers_01_two_sum_sorted`) — Easy — *O(N) Time, O(1) Space*
2. **Maximum Subarray (Kadane's Algorithm)** (`twopointers_02_maximum_subarray`) — Medium — *O(N) Time, O(1) Space*

---

## 🛠️ CLI Subcommands

In addition to the interactive TUI, `algo` provides standalone CLI commands suitable for terminal shortcuts and CI pipelines:

```bash
# List all tracks, exercises, and completion badges
algo list

# Verify an exercise without opening the TUI
algo verify sorting_01_bubble_sort

# Show hints for an exercise
algo hint sorting_01_bubble_sort

# Display the reference solution for an exercise
algo solution sorting_01_bubble_sort

# Extract embedded exercises into another folder
algo init ./my-algo-workspace
```

---

## 🧩 Adding New Exercises

To add a new exercise to the platform, create a directory under `exercises/<track_name>/<exercise_name>` containing:

1. **`exercise.json`:**
   ```json
   {
     "id": "graph_01_bfs",
     "name": "Breadth-First Search",
     "track": "Graphs",
     "track_order": 6,
     "order": 1,
     "difficulty": "Medium",
     "description": "Traverse an unweighted graph using BFS...",
     "complexity": { "time": "O(V + E)", "space": "O(V)" },
     "forbidden": ["container/list"],
     "entry_file": "bfs.go",
     "test_file": "bfs_test.go",
     "hints": ["Use a queue to track nodes to visit..."]
   }
   ```
2. **Starter code** (e.g. `bfs.go`) with stub function.
3. **Unit tests & benchmarks** (e.g. `bfs_test.go`).
4. **Reference solution** under `solution/solution.go`.

The platform auto-discovers all exercises dynamically!

---

## 🧪 Running Internal Tests

```bash
go test -v ./internal/...
```
