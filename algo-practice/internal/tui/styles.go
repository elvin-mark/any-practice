package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#7D56F4")
	colorSecondary = lipgloss.Color("#5A32A3")
	colorSuccess   = lipgloss.Color("#04B575")
	colorError     = lipgloss.Color("#FF4040")
	colorWarning   = lipgloss.Color("#FFAE00")
	colorMuted     = lipgloss.Color("#666666")
	colorSubtle    = lipgloss.Color("#888888")
	colorBgDark    = lipgloss.Color("#1A1A24")
	colorTextLight = lipgloss.Color("#FFFFFF")
	colorHighlight = lipgloss.Color("#00D7D7")

	// Header Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorTextLight).
			Background(colorPrimary).
			Padding(0, 1)

	trackBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Background(lipgloss.Color("#2A2438")).
			Padding(0, 1).
			MarginLeft(1)

	progressBadgeStyle = lipgloss.NewStyle().
				Foreground(colorHighlight).
				Bold(true).
				MarginLeft(1)

	// Difficulty Badges
	diffEasyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess).
			Background(lipgloss.Color("#1B3320")).
			Padding(0, 1)

	diffMedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning).
			Background(lipgloss.Color("#332A15")).
			Padding(0, 1)

	diffHardStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorError).
			Background(lipgloss.Color("#331818")).
			Padding(0, 1)

	// Panel Box Styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3854")).
			Padding(0, 1).
			MarginBottom(1)

	successCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorSuccess).
				Padding(0, 1).
				MarginBottom(1)

	errorCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError).
			Padding(0, 1).
			MarginBottom(1)

	// Text formatting
	boldText   = lipgloss.NewStyle().Bold(true).Foreground(colorTextLight)
	subtleText = lipgloss.NewStyle().Foreground(colorSubtle)
	mutedText  = lipgloss.NewStyle().Foreground(colorMuted)

	sectionHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHighlight).
			MarginTop(1)

	// Footer / Keybinding Style
	keyBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	keyDescStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	footerStyle = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#333333")).
			PaddingTop(0)
)
