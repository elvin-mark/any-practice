package ui

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

// Legacy color functions for backward compatibility
var (
	CyanBold    = color.New(color.FgCyan, color.Bold).SprintFunc()
	GreenBold   = color.New(color.FgGreen, color.Bold).SprintFunc()
	RedBold     = color.New(color.FgRed, color.Bold).SprintFunc()
	YellowBold  = color.New(color.FgYellow, color.Bold).SprintFunc()
	MagentaBold = color.New(color.FgMagenta, color.Bold).SprintFunc()
	BlueBold    = color.New(color.FgBlue, color.Bold).SprintFunc()
	WhiteBold   = color.New(color.FgWhite, color.Bold).SprintFunc()
	Gray        = color.New(color.FgHiBlack).SprintFunc()
	Dim         = color.New(color.Faint).SprintFunc()

	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
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

// PrintBanner prints the unified CKA header banner
func PrintBanner(args ...string) {
	toolTitle := "CKA-PRACTICE"
	subtitle := "Certified Kubernetes Administrator Interactive Exam & Practice Labs"
	if len(args) >= 1 {
		toolTitle = args[0]
	}
	if len(args) >= 2 {
		subtitle = args[1]
	}

	badge := TitleBadgeStyle.Render(" " + toolTitle + " ")
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
	case strings.Contains(domain, "Architecture"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#B388FF")).Bold(true).Render("[ARCH/CFG]")
	case strings.Contains(domain, "Workloads"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#448AFF")).Bold(true).Render("[WORKLOAD]")
	case strings.Contains(domain, "Services"):
		return lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("[NETWORK ]")
	case strings.Contains(domain, "Storage"):
		return lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("[STORAGE ]")
	case strings.Contains(domain, "Troubleshooting"):
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

// FormatStatus returns a colored completion status indicator
func FormatStatus(completed bool) string {
	if completed {
		return StatusCompleted
	}
	return StatusPending
}

// PrintBox wraps text in a clean Unicode card box
func PrintBox(title, content string) {
	boxContent := fmt.Sprintf("%s\n\n%s", BoldText.Render(title), content)
	fmt.Println(CardStyle.Width(76).Render(boxContent))
}

// TableColumn configures column header and min width
type TableColumn struct {
	Title string
	Width int
}

// RenderTable prints a gorgeous Unicode-bordered table with ANSI-aware cell padding
func RenderTable(columns []TableColumn, rows [][]string) {
	if len(columns) == 0 {
		return
	}

	// Compute effective column widths
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

	// Build border lines
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
