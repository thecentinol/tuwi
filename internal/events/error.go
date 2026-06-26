package events

import tea "charm.land/bubbletea/v2"

type ShowErrorMsg struct{ Err error }
type DismissErrorMsg struct{}

func ShowError(err error) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorMsg{Err: err}
	}
}

func DismissError() tea.Cmd {
	return func() tea.Msg {
		return DismissErrorMsg{}
	}
}
