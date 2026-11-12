# sql-practice (sqlp)

> Terminal-native SQL mastery, query optimization, recursive CTEs, and schema indexing labs.

`sqlp` provides automated hands-on practice labs for modern SQL using fast, pure-Go embedded SQLite databases:
- **Core Queries & Aggregations**: Grouping, filtering, sorting, multi-table joins.
- **Advanced Joins**: Left joins, anti-joins for finding missing or unpurchased relations.
- **Window Functions**: Running totals, cumulative metrics (`SUM() OVER`), and rank analytics (`DENSE_RANK()`).
- **Common Table Expressions**: Recursive CTEs for hierarchy depth, graph traversal, and reporting trees.
- **Index Optimization**: Creating composite indexes to eliminate slow full-table scans verified by `EXPLAIN QUERY PLAN`.

Each lab seeds an isolated SQLite database at `/tmp/sqlp-labs/<scenario-id>/practice.db` and generates a `solution.sql` workspace file.

---

## Quickstart

```bash
# Build binary
make build

# Launch interactive split-pane TUI
./bin/sqlp

# List all available labs
./bin/sqlp list

# Start a specific lab (seeds database and creates solution.sql)
./bin/sqlp start query-01-top-customers

# Run an ad-hoc SQL query and view results in a Unicode table:
./bin/sqlp query query-01-top-customers "SELECT * FROM customers LIMIT 5;"

# Write your solution query in /tmp/sqlp-labs/query-01-top-customers/solution.sql
# Then verify:
./bin/sqlp verify query-01-top-customers

# Or verify with an inline query:
./bin/sqlp verify query-01-top-customers "SELECT c.id, c.name, SUM(o.amount) AS total FROM customers c JOIN orders o ON c.id=o.customer_id GROUP BY c.id ORDER BY total DESC LIMIT 3;"

# Inspect database schema and tables
./bin/sqlp tables query-01-top-customers

# View progress breakdown
./bin/sqlp progress
```

---

## Interactive TUI Controls

| Key | Action |
| --- | --- |
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `Enter` / `s` | Initialize database & start lab |
| `v` | Verify solution in `solution.sql` |
| `r` | Reset database to initial seeded state |
| `h` | Reveal next progressive hint |
| `q` / `Ctrl+C` | Quit |

---

## Adding New Scenarios

Add new YAML scenario definitions into `scenarios/`:

```yaml
id: my-query-lab
title: "Lab Name"
category: Queries
difficulty: Easy
description: Problem statement...
goal: Expected query output...
hints:
  - First hint
setup: |
  CREATE TABLE items (id INT, price REAL);
  INSERT INTO items VALUES (1, 10.0), (2, 25.5);
solution: |
  SELECT id, price FROM items WHERE price > 15;
```
