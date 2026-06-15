package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/wifi"
)

type (
	WifiModel struct {
		client *dbus.Client

		savedList     WifiListModel
		availableList WifiListModel

		scanning bool
	}

	savedWifiMsg     []wifi.AccessPoint
	availableWifiMsg []wifi.AccessPoint
)

func (w WifiModel) Init() tea.Cmd {
	return tea.Batch(
		fetchAvailableWifiNetworks(w.client),
	)
}

func (w WifiModel) Update(msg tea.Msg) (WifiModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case savedWifiMsg:
		w.savedList.networks = msg
	case availableWifiMsg:
		w.scanning = false
		w.availableList.networks = msg
		return w, fetchSavedNetworks(w.client, []wifi.AccessPoint(msg))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			w.scanning = true
			cmds = append(cmds, fetchAvailableWifiNetworks(w.client))
		}
	}

	return w, tea.Batch(cmds...)
}

func (w WifiModel) View() tea.View {
	return tea.NewView(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			w.savedList.View().Content,
			w.availableList.View().Content,
		),
	)
}

func fetchSavedNetworks(c *dbus.Client, available []wifi.AccessPoint) tea.Cmd {
	return func() tea.Msg {
		saved, err := wifi.GetSavedNetworks(c, available)
		if err != nil {
			return err
		}
		return savedWifiMsg(saved)
	}
}

func fetchAvailableWifiNetworks(c *dbus.Client) tea.Cmd {
	var w WifiModel
	return func() tea.Msg {
		w.scanning = true
		networks, err := wifi.GetAvailableNetworks(c)
		if err != nil {
			return err
		}
		return availableWifiMsg(networks)
	}
}
