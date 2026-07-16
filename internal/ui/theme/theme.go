package theme

import (
	"charm.land/lipgloss/v2"
)

var (
	Base        = lipgloss.Color("238")
	Focused     = lipgloss.Green
	FocusedLine = lipgloss.NewStyle().
			Background(Focused).
			Foreground(Base)

	// focusedBorderStyle = lipgloss.NewStyle().
	// 			Border(lipgloss.RoundedBorder()).
	// 			BorderForeground(lipgloss.Color("238"))
	//
	// blurredBorderStyle = lipgloss.NewStyle().
	// 			Border(lipgloss.HiddenBorder())
)
