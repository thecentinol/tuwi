package events

import (
	tea "charm.land/bubbletea/v2"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

type ShowPasswordModalMsg struct {
	Network *nm.AccessPoint
}

type PasswordResultMsg struct {
	Cancelled bool
	Password  string
}

type OpenHelpMenuMsg struct{}

type CloseHelpMenuMsg struct{}

// returns OpenHelpMenuMsg{}
func OpenHelpMenu() tea.Cmd {
	return func() tea.Msg {
		return OpenHelpMenuMsg{}
	}
}

// returns CloseHelpMenuMsg{}
func CloseHelpMenu() tea.Cmd {
	return func() tea.Msg {
		return CloseHelpMenuMsg{}
	}
}
