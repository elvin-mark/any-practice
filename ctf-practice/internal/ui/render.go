package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Design System Color Tokens
var (
	ColorPrimary    = lipgloss.Color("#7D56F4")
	ColorSecondary  = lipgloss.Color("#5A32A3")
	ColorAccent     = lipgloss.Color("#00D7D7")
	ColorSuccess    = lipgloss.Color("#04B575")
	ColorWarning    = lipgloss.Color("#FFAE00")
	ColorError      = lipgloss.Color("#FF4040")
	ColorDarkBg     = lipgloss.Color("#1A1A24")
	ColorCardBorder = lipgloss.Color("#3C3854")
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

	PointsBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAccent).
				Background(lipgloss.Color("#162938")).
				Padding(0, 1)

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

// FormatDifficulty returns a colored badge for difficulty levels
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

// FormatStatus returns a colored completion status indicator
func FormatStatus(completed bool) string {
	if completed {
		return StatusCompleted
	}
	return StatusPending
}

// FormatPoints returns a formatted points indicator
func FormatPoints(points int) string {
	return PointsBadgeStyle.Render(fmt.Sprintf("%d pts", points))
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

// PrintBanner renders the unified tool banner
func PrintBanner(toolTitle, subtitle string) {
	badge := TitleBadgeStyle.Render(" " + toolTitle + " ")
	sub := SubtleText.Render(subtitle)
	line := SubtleText.Render(strings.Repeat("─", 74))

	fmt.Printf("\n%s\n  %s  %s\n%s\n", line, badge, sub, line)
}

// PrintHeader prints a section divider
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

// PrintFailure prints an error/failure message
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

// PrintBox prints text inside a rounded Unicode card box
func PrintBox(title, content string) {
	boxContent := fmt.Sprintf("%s\n\n%s", BoldText.Render(title), content)
	fmt.Println(CardStyle.Width(76).Render(boxContent))
}

// TableColumn represents a table column definition
type TableColumn struct {
	Title string
	Width int
}

// RenderTable prints a Unicode-bordered table matching the unified design system
func RenderTable(columns []TableColumn, rows [][]string) {
	if len(columns) == 0 {
		return
	}

	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = lipgloss.Width(col.Title)
		if col.Width > widths[i] {
			widths[i] = col.Width
		}
	}

	for _, row := range rows {
		for i := 0; i < len(columns) && i < len(row); i++ {
			w := lipgloss.Width(row[i])
			if w > widths[i] {
				widths[i] = w
			}
		}
	}

	// Border builder
	var topBorder, sepBorder, botBorder strings.Builder
	topBorder.WriteString("┌")
	sepBorder.WriteString("├")
	botBorder.WriteString("└")

	for i, w := range widths {
		pad := strings.Repeat("─", w+2)
		topBorder.WriteString(pad)
		sepBorder.WriteString(pad)
		botBorder.WriteString(pad)
		if i < len(widths)-1 {
			topBorder.WriteString("┬")
			sepBorder.WriteString("┼")
			botBorder.WriteString("┴")
		} else {
			topBorder.WriteString("┐")
			sepBorder.WriteString("┤")
			botBorder.WriteString("┘")
		}
	}

	borderStyle := lipgloss.NewStyle().Foreground(ColorCardBorder)

	// Print Top Border
	fmt.Println(borderStyle.Render(topBorder.String()))

	// Print Header
	var headerLine strings.Builder
	headerLine.WriteString(borderStyle.Render("│"))
	for i, col := range columns {
		padWidth := widths[i] - lipgloss.Width(col.Title)
		headerLine.WriteString(" " + BoldText.Render(col.Title) + strings.Repeat(" ", padWidth) + " " + borderStyle.Render("│"))
	}
	fmt.Println(headerLine.String())

	// Print Separator
	fmt.Println(borderStyle.Render(sepBorder.String()))

	// Print Rows
	for _, row := range rows {
		var rowLine strings.Builder
		rowLine.WriteString(borderStyle.Render("│"))
		for i := 0; i < len(columns); i++ {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			padWidth := widths[i] - lipgloss.Width(val)
			if padWidth < 0 {
				padWidth = 0
			}
			rowLine.WriteString(" " + val + strings.Repeat(" ", padWidth) + " " + borderStyle.Render("│"))
		}
		fmt.Println(rowLine.String())
	}

	// Print Bottom Border
	fmt.Println(borderStyle.Render(botBorder.String()))
}
