package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/models"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/wifi"
)

type (
	WifiModel struct {
		client *dbus.Client

		savedList     WifiListModel
		availableList WifiListModel

		scanning bool
		err      error
	}

	savedWifiMsg     []models.AccessPoint
	availableWifiMsg []models.AccessPoint
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

		rows := comp.AccessPointsToRows(msg)
		w.savedList.table.SetRows(rows)
	case availableWifiMsg:
		w.scanning = false
		w.availableList.networks = msg

		rows := comp.AccessPointsToRows(msg)
		w.availableList.table.SetRows(rows)
		return w, fetchSavedNetworks(w.client, []models.AccessPoint(msg))

	case error:
		w.scanning = false
		w.err = msg
	}

	var cmd tea.Cmd
	if w.savedList.focused {
		w.savedList, cmd = w.savedList.Update(msg)
	} else {
		w.availableList, cmd = w.availableList.Update(msg)
	}
	cmds = append(cmds, cmd)

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

func (w WifiModel) HelpView() []key.Binding {
	if w.savedList.focused {
		return w.savedList.HelpView()
	}
	return w.availableList.HelpView()
}

func fetchSavedNetworks(c *dbus.Client, available []models.AccessPoint) tea.Cmd {
	return func() tea.Msg {
		saved, err := wifi.GetSavedNetworks(c, available)
		if err != nil {
			return err
		}
		return savedWifiMsg(saved)
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
