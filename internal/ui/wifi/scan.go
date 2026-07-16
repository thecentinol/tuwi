package wifi

import (
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

func handleScan(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := wifidomain.GetAvailableNetworks(c)
		if err != nil {
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}
