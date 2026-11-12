package ui

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Design System Color Tokens
var (
	ColorPrimary    = lipgloss.Color("#7D56F4") // Electric Violet
	ColorSecondary  = lipgloss.Color("#5A32A3") // Deep Purple
	ColorAccent     = lipgloss.Color("#00D7D7") // Cyan
	ColorSuccess    = lipgloss.Color("#04B575") // Emerald Green
	ColorWarning    = lipgloss.Color("#FFAE00") // Amber Gold
	ColorError      = lipgloss.Color("#FF4040") // Coral Red
	ColorDarkBg     = lipgloss.Color("#1A1A24")
	ColorCardBorder = lipgloss.Color("#3C3854") // Muted Purple-Slate
	ColorActiveBg   = lipgloss.Color("#353366")
	ColorText       = lipgloss.Color("#FAFAFA")
	ColorSubtle     = lipgloss.Color("#888888")
	ColorMuted      = lipgloss.Color("#666666")
)

// Badges and Text Styles
var (
	TitleBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	CategoryBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Background(lipgloss.Color("#2A2438")).
				Padding(0, 1)

	ProgressBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAccent)

	DiffEasyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess).
			Background(lipgloss.Color("#1B3320")).
			Padding(0, 1)

	DiffMedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning).
			Background(lipgloss.Color("#332A15")).
			Padding(0, 1)

	DiffHardStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError).
			Background(lipgloss.Color("#331818")).
			Padding(0, 1)

	StatusCompleted = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess).
			Render("[✓]")

	StatusPending = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Render("[ ]")

	StatusActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Render("[▶]")

	BoldText   = lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	SubtleText = lipgloss.NewStyle().Foreground(ColorSubtle)
	MutedText  = lipgloss.NewStyle().Foreground(ColorMuted)
	CyanText   = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	YellowText = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	GreenText  = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	RedText    = lipgloss.NewStyle().Foreground(ColorError).Bold(true)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCardBorder).
			Padding(0, 1)

	SuccessCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSuccess).
				Padding(0, 1)

	ErrorCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Padding(0, 1)

	KeyBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle)
)

// VisibleLen calculates rune length without ANSI escape codes
func VisibleLen(s string) int {
	clean := ansiRegex.ReplaceAllString(s, "")
	return utf8.RuneCountInString(clean)
}

// PadRight pads string with spaces based on visible width
func PadRight(s string, width int) string {
	vLen := VisibleLen(s)
	if vLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vLen)
}

// PadLeft pads string on the left
func PadLeft(s string, width int) string {
	vLen := VisibleLen(s)
	if vLen >= width {
		return s
	}
	return strings.Repeat(" ", width-vLen) + s
}

// RenderKeyHelp formats a slice of key-action pairs into a unified footer
func RenderKeyHelp(keys [][2]string) string {
	var parts []string
	for _, k := range keys {
		part := fmt.Sprintf("%s %s", KeyBadgeStyle.Render(k[0]), KeyDescStyle.Render(k[1]))
		parts = append(parts, part)
	}
	return strings.Join(parts, SubtleText.Render(" • "))
}

// PrintBanner prints the unified Docker header banner
func PrintBanner(args ...string) {
	toolTitle := "DOCKER-PRACTICE"
	subtitle := "Hands-on Docker Containerization, Production Architecture & Security"
	if len(args) >= 1 && args[0] != "" {
		toolTitle = args[0]
	}
	if len(args) >= 2 && args[1] != "" {
		subtitle = args[1]
	}

	badge := TitleBadgeStyle.Render(" 🐳 " + toolTitle + " ")
	sub := SubtleText.Render(subtitle)
	line := SubtleText.Render(strings.Repeat("─", 74))

	fmt.Printf("\n%s\n  %s  %s\n%s\n", line, badge, sub, line)
}

// PrintHeader prints a section header with border
func PrintHeader(title string) {
	line := SubtleText.Render(strings.Repeat("─", 74))
	fmt.Printf("\n%s\n  %s\n%s\n", line, CyanText.Render(title), line)
}

// PrintSubHeader prints a sub-section
func PrintSubHeader(title string) {
	fmt.Printf("\n%s %s\n", YellowText.Render("▶"), BoldText.Render(title))
}

// PrintSuccess prints a success message
func PrintSuccess(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", GreenText.Render("✔ PASS:"), msg)
}

// PrintFailure prints a failure message
func PrintFailure(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", RedText.Render("✖ FAIL:"), msg)
}

// PrintInfo prints an informational message
func PrintInfo(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", CyanText.Render("ℹ INFO:"), msg)
}

// PrintWarning prints a warning message
func PrintWarning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", YellowText.Render("⚠ WARN:"), msg)
}

// PrintHint prints a hint message
func PrintHint(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", YellowText.Render("💡 HINT:"), msg)
}

// FormatDomain returns a colored badge for domain names
func FormatDomain(domain string) string {
	switch {
	case strings.Contains(strings.ToLower(domain), "image") || strings.Contains(strings.ToLower(domain), "dockerfile"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#B388FF")).Bold(true).Render("[IMAGES  ]")
	case strings.Contains(strings.ToLower(domain), "lifecycle") || strings.Contains(strings.ToLower(domain), "resource"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#448AFF")).Bold(true).Render("[LIFECYCL]")
	case strings.Contains(strings.ToLower(domain), "network"):
		return lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("[NETWORK ]")
	case strings.Contains(strings.ToLower(domain), "storage") || strings.Contains(strings.ToLower(domain), "volume"):
		return lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("[STORAGE ]")
	case strings.Contains(strings.ToLower(domain), "compose"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF80AB")).Bold(true).Render("[COMPOSE ]")
	case strings.Contains(strings.ToLower(domain), "security"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5252")).Bold(true).Render("[SECURITY]")
	case strings.Contains(strings.ToLower(domain), "trouble"):
		return lipgloss.NewStyle().Foreground(ColorError).Bold(true).Render("[TROUBLE ]")
	default:
		return SubtleText.Render("[GENERAL ]")
	}
}

// FormatDifficulty returns a colored difficulty badge
func FormatDifficulty(diff string) string {
	switch strings.ToLower(strings.TrimSpace(diff)) {
	case "easy":
		return DiffEasyStyle.Render("Easy")
	case "medium", "med":
		return DiffMedStyle.Render("Medium")
	case "hard":
		return DiffHardStyle.Render("Hard")
	default:
		return DiffMedStyle.Render(diff)
	}
}

// TableColumn represents a column header and width
type TableColumn struct {
	Title string
	Width int
}

// RenderTable renders a stylized Unicode box table
func RenderTable(columns []TableColumn, rows [][]string) string {
	var b strings.Builder

	// Top border: ┌───┬───┐
	b.WriteString("┌")
	for i, col := range columns {
		b.WriteString(strings.Repeat("─", col.Width+2))
		if i < len(columns)-1 {
			b.WriteString("┬")
		}
	}
	b.WriteString("┐\n")

	// Header row: │ Title │
	b.WriteString("│")
	for _, col := range columns {
		b.WriteString(" " + BoldText.Render(PadRight(col.Title, col.Width)) + " │")
	}
	b.WriteString("\n")

	// Header separator: ├───┼───┤
	b.WriteString("├")
	for i, col := range columns {
		b.WriteString(strings.Repeat("─", col.Width+2))
		if i < len(columns)-1 {
			b.WriteString("┼")
		}
	}
	b.WriteString("┤\n")

	// Data rows: │ Cell │
	for _, row := range rows {
		b.WriteString("│")
		for j, col := range columns {
			val := ""
			if j < len(row) {
				val = row[j]
			}
			b.WriteString(" " + PadRight(val, col.Width) + " │")
		}
		b.WriteString("\n")
	}

	// Bottom border: └───┴───┘
	b.WriteString("└")
	for i, col := range columns {
		b.WriteString(strings.Repeat("─", col.Width+2))
		if i < len(columns)-1 {
			b.WriteString("┴")
		}
	}
	b.WriteString("┘\n")

	return b.String()
}
