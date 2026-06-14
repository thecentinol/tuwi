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

		scanning bool // TODO: implement me!
	}

	savedWifiMsg     []wifi.AccessPoint
	availableWifiMsg []wifi.AccessPoint
)

func (w WifiModel) Init() tea.Cmd {
	return tea.Batch(
		fetchAvailableWifiNetworks(w.client),
		fetchSavedNetworks(w.client),
	)
}

func (w WifiModel) Update(msg tea.Msg) (WifiModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case savedWifiMsg:
		w.savedList.networks = msg
	case availableWifiMsg:
		w.availableList.networks = msg

	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
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

func fetchSavedNetworks(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := wifi.GetSavedNetworks(c)
		if err != nil {
			return err
		}
		return savedWifiMsg(networks)
	}
}

func fetchAvailableWifiNetworks(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := wifi.GetAvailableNetworks(c)
		if err != nil {
			return err
		}
		return availableWifiMsg(networks)
	}
}
