package ui

import (
	"charm.land/lipgloss/v2"
)

var (
	base        = lipgloss.Color("238")
	focused     = lipgloss.Green
	focusedLine = lipgloss.NewStyle().
			Background(focused).
			Foreground(base)

	// focusedBorderStyle = lipgloss.NewStyle().
	// 			Border(lipgloss.RoundedBorder()).
	// 			BorderForeground(lipgloss.Color("238"))
	//
	// blurredBorderStyle = lipgloss.NewStyle().
	// 			Border(lipgloss.HiddenBorder())
)
