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

type OpenKeybindsModalMsg struct{}

type CloseKeybindsModalMsg struct{}

// returns OpenKeybindsModalMsg{}
func OpenKeybindsModal() tea.Cmd {
	return func() tea.Msg {
		return OpenKeybindsModalMsg{}
	}
}

// returns CloseKeybindsModalMsg{}
func CloseKeybindsModal() tea.Cmd {
	return func() tea.Msg {
		return CloseKeybindsModalMsg{}
	}
}
