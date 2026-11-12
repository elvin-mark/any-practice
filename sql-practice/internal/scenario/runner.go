package scenario

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// LabDir returns the directory containing the SQLite database and solution.sql for a scenario.
func LabDir(scenarioID string) string {
	return filepath.Join(os.TempDir(), "sqlp-labs", scenarioID)
}

// DBPath returns the path to the SQLite database file for a scenario.
func DBPath(scenarioID string) string {
	return filepath.Join(LabDir(scenarioID), "practice.db")
}

// SolutionFilePath returns the path to the user's solution.sql file.
func SolutionFilePath(scenarioID string) string {
	return filepath.Join(LabDir(scenarioID), "solution.sql")
}

// SetupLab initializes the SQLite database with DDL and test data.
func SetupLab(sc *Scenario) (string, error) {
	dir := LabDir(sc.ID)
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating lab directory: %w", err)
	}

	dbPath := DBPath(sc.ID)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", fmt.Errorf("opening sqlite database: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(sc.Setup); err != nil {
		return "", fmt.Errorf("executing scenario setup DDL/DML: %w", err)
	}

	// Create starter solution.sql file
	solHeader := fmt.Sprintf("-- Scenario: %s\n-- Goal: %s\n-- Write your SQL query below and run: sqlp verify %s\n\n",
		sc.Title, sc.Goal, sc.ID)
	_ = os.WriteFile(SolutionFilePath(sc.ID), []byte(solHeader), 0644)

	return dbPath, nil
}

// ExecuteQuery runs an arbitrary query against the scenario's SQLite database and returns columns and formatted rows.
func ExecuteQuery(dbPath, query string) ([]string, [][]string, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("opening database %s: %w", dbPath, err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("fetching columns: %w", err)
	}

	var results [][]string
	for rows.Next() {
		colValues := make([]interface{}, len(cols))
		colPointers := make([]interface{}, len(cols))
		for i := range colValues {
			colPointers[i] = &colValues[i]
		}

		if err := rows.Scan(colPointers...); err != nil {
			return nil, nil, fmt.Errorf("scanning row: %w", err)
		}

		rowStrs := make([]string, len(cols))
		for i, val := range colValues {
			if val == nil {
				rowStrs[i] = "NULL"
			} else {
				rowStrs[i] = fmt.Sprintf("%v", val)
			}
		}
		results = append(results, rowStrs)
	}

	return cols, results, rows.Err()
}

// VerifyLab validates the user's query against the canonical solution.
func VerifyLab(sc *Scenario, userSQL string) (bool, string, error) {
	dbPath := DBPath(sc.ID)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false, fmt.Sprintf("Lab database does not exist. Run 'sqlp start %s' first.", sc.ID), nil
	}

	// If no inline query was provided, read from solution.sql
	if strings.TrimSpace(userSQL) == "" {
		solFile := SolutionFilePath(sc.ID)
		data, err := os.ReadFile(solFile)
		if err != nil {
			return false, fmt.Sprintf("Could not read solution file %s: %v", solFile, err), nil
		}
		userSQL = string(data)
	}

	cleanUserSQL := cleanQuery(userSQL)
	if cleanUserSQL == "" {
		return false, "No SQL query found. Write your query in " + SolutionFilePath(sc.ID) + " or pass via argument.", nil
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return false, "", fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Special verification mode: Index Optimization / Query Plan Check
	if sc.VerifyType == "plan_match" {
		// First execute user's statement (which may be a CREATE INDEX statement or query)
		if _, err := db.Exec(cleanUserSQL); err != nil {
			return false, fmt.Sprintf("SQL Execution Error: %v", err), nil
		}

		// Now explain the reference query
		explainQuery := "EXPLAIN QUERY PLAN " + sc.Solution
		_, rows, err := ExecuteQuery(dbPath, explainQuery)
		if err != nil {
			return false, fmt.Sprintf("Error explaining plan: %v", err), nil
		}

		planDump := ""
		for _, r := range rows {
			planDump += strings.Join(r, " | ") + "\n"
		}

		targetPlan := strings.ToUpper(strings.TrimSpace(sc.VerifyPlan))
		if strings.Contains(strings.ToUpper(planDump), targetPlan) {
			return true, fmt.Sprintf("Query plan successfully matches expectation!\n%s", planDump), nil
		}
		return false, fmt.Sprintf("Query plan does not use expected index (%s).\nCurrent Plan:\n%s", sc.VerifyPlan, planDump), nil
	}

	// Standard result-matching verification
	_, expectedRows, err := ExecuteQuery(dbPath, sc.Solution)
	if err != nil {
		return false, "", fmt.Errorf("executing canonical solution: %w", err)
	}

	userCols, userRows, err := ExecuteQuery(dbPath, cleanUserSQL)
	if err != nil {
		return false, fmt.Sprintf("Query syntax or execution error: %v", err), nil
	}

	// 1. Check Row Count
	if len(userRows) != len(expectedRows) {
		return false, fmt.Sprintf("Row count mismatch: expected %d rows, got %d rows.", len(expectedRows), len(userRows)), nil
	}

	// 2. Check Column Count
	if len(expectedRows) > 0 {
		expColsCount := len(expectedRows[0])
		if len(userCols) != expColsCount {
			return false, fmt.Sprintf("Column count mismatch: expected %d columns, got %d columns.", expColsCount, len(userCols)), nil
		}
	}

	// 3. Check Row Values
	for i := range expectedRows {
		for j := range expectedRows[i] {
			if strings.TrimSpace(userRows[i][j]) != strings.TrimSpace(expectedRows[i][j]) {
				return false, fmt.Sprintf("Value mismatch at row %d, col %d ('%s'): expected '%s', got '%s'.",
					i+1, j+1, userCols[j], expectedRows[i][j], userRows[i][j]), nil
			}
		}
	}

	return true, fmt.Sprintf("All %d rows and %d columns matched expected output perfectly!", len(userRows), len(userCols)), nil
}

// cleanQuery strips SQL comments and extra whitespace.
func cleanQuery(q string) string {
	lines := strings.Split(q, "\n")
	var kept []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		kept = append(kept, l)
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// GetTables returns a list of user tables in the database.
func GetTables(dbPath string) ([]string, error) {
	_, rows, err := ExecuteQuery(dbPath, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name;")
	if err != nil {
		return nil, err
	}
	var tables []string
	for _, r := range rows {
		if len(r) > 0 {
			tables = append(tables, r[0])
		}
	}
	return tables, nil
}
