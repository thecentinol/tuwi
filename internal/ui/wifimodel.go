package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/models"
)

type (
	WifiListModel struct {
		width    int
		height   int
		focused  bool
		networks []models.AccessPoint
		cursor   int
	}

	savedWifiMsg     []models.AccessPoint
	availableWifiMsg []models.AccessPoint
)

func (w WifiListModel) Update(msg tea.Msg) (WifiListModel, tea.Cmd) {
	return w, nil
}

func (w WifiListModel) View() tea.View {
	content := ""
	borderColor := base
	if w.focused {
		borderColor = focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(w.width).
		Height(w.height)

	for i, v := range w.networks {
		if i == w.cursor {
			content += "> " + v.SSID + "\n"
		} else {
			content += "  " + v.SSID + "\n"
		}
	}

	return tea.NewView(style.Render(content))
}

func fetchSavedNetworks(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := dbus.GetSavedNetworks(c)
		if err != nil {
			return err
		}
		return savedWifiMsg(networks)
	}
}

func fetchAvailableWifiNetworks(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := dbus.GetNetworks(c)
		if err != nil {
			return err
		}
		return availableWifiMsg(networks)
	}
}
