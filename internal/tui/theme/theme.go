package theme

import (
	"charm.land/lipgloss/v2"
)

// Catppuccin Mocha-inspired palette
var (
	ColorBase     = lipgloss.Color("#1e1e2e")
	ColorSurface0 = lipgloss.Color("#313244")
	ColorSurface1 = lipgloss.Color("#45475a")
	ColorSurface2 = lipgloss.Color("#585b70")
	ColorOverlay0 = lipgloss.Color("#6c7086")
	ColorOverlay1 = lipgloss.Color("#7f849c")
	ColorText     = lipgloss.Color("#cdd6f4")
	ColorSubtext0 = lipgloss.Color("#a6adc8")
	ColorSubtext1 = lipgloss.Color("#bac2de")
	ColorLavender = lipgloss.Color("#b4befe")
	ColorBlue     = lipgloss.Color("#89b4fa")
	ColorSapphire = lipgloss.Color("#74c7ec")
	ColorGreen    = lipgloss.Color("#a6e3a1")
	ColorYellow   = lipgloss.Color("#f9e2af")
	ColorPeach    = lipgloss.Color("#fab387")
	ColorRed      = lipgloss.Color("#f38ba8")
	ColorMauve    = lipgloss.Color("#cba6f7")
	ColorPink     = lipgloss.Color("#f5c2e7")
	ColorTeal     = lipgloss.Color("#94e2d5")
	ColorFlamingo = lipgloss.Color("#f2cdcd")
)

// App-level styles
var (
	AppStyle = lipgloss.NewStyle().
			Background(ColorBase)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorLavender).
			Padding(0, 1)

	HeaderTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorMauve)

	HeaderFilterStyle = lipgloss.NewStyle().
				Foreground(ColorSubtext0).
				Padding(0, 1)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorLavender).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorSurface1).
				Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
				Background(ColorSurface1).
				Foreground(ColorText)

	NormalRowStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	DimRowStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay1)

	UnreadIndicator = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	ReadIndicator = lipgloss.NewStyle().
			Foreground(ColorSurface2)

	RepoStyle = lipgloss.NewStyle().
			Foreground(ColorSapphire)

	NumberStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	TimeStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay0)

	// Preview panel
	PreviewBorderStyle = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorSurface1).
				Padding(1, 2)

	PreviewTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorLavender).
				MarginBottom(1)

	PreviewMetaStyle = lipgloss.NewStyle().
				Foreground(ColorSubtext0)

	PreviewLabelStyle = lipgloss.NewStyle().
				Foreground(ColorBase).
				Background(ColorMauve).
				Padding(0, 1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext0).
			Background(ColorSurface0).
			Padding(0, 1)

	StatusBarKeyStyle = lipgloss.NewStyle().
				Foreground(ColorLavender).
				Bold(true)

	StatusBarValueStyle = lipgloss.NewStyle().
				Foreground(ColorText)

	// Filter / search
	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(ColorYellow).
				Bold(true)

	// Help
	HelpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMauve).
			MarginBottom(1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true).
			Width(14)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext0)

	// Notifications / feedback
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ColorMauve)
)
