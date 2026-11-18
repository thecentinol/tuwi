package wifi

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/service"
)

func handleScan(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := service.GetAvailableNetworks(c)
		if err != nil {
			slog.Error("Failed network scan", "error", err)
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}
