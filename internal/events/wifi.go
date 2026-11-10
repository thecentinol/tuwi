package events

import (
	tea "charm.land/bubbletea/v2"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

type (
	// initial saved connections have loaded
	SavedWifiMsg struct{ Networks []nm.SavedConnection }
	// initial available access points have loaded
	AvailableWifiMsg struct{ Networks []nm.AccessPoint }

	// used when the saved or available wifi list has changed
	// so the info pane must be updated
	WifiListChangedMsg struct{}
)

type WifiConnectReqMsg struct {
	Network  *nm.AccessPoint
	Password string
}

func WifiListChanged() tea.Cmd {
	return func() tea.Msg {
		return WifiListChangedMsg{}
	}
}
