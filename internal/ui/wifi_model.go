package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	godbus "github.com/godbus/dbus/v5"
	"slices"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/wifi"
)

type WifiModel struct {
	client *dbus.Client

	savedList     WifiListModel
	availableList WifiListModel

	scanning bool
	err      error
}

func (w WifiModel) Init() tea.Cmd {
	return tea.Batch(
		fetchAvailableWifiNetworks(w.client),
	)
}

func (w WifiModel) Update(msg tea.Msg) (WifiModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case events.SavedWifiMsg:
		w.savedList.networks = msg.Networks
		w.savedList.table.SetRows(AccessPointsToRows(msg.Networks))

		savedBSSIDs := make(map[string]bool)
		for _, saved := range w.savedList.networks {
			if saved.BSSID != "" {
				savedBSSIDs[saved.BSSID] = true
			}
		}

		result := slices.DeleteFunc(w.availableList.networks, func(n nm.AccessPoint) bool {
			return savedBSSIDs[n.BSSID]
		})
		w.availableList.networks = result
		w.availableList.table.SetRows(AccessPointsToRows(result))

	case events.AvailableWifiMsg:
		w.scanning = false
		w.availableList.networks = msg.Networks

		rows := AccessPointsToRows(msg.Networks)
		w.availableList.table.SetRows(rows)
		cmds = append(cmds, fetchSavedNetworks(w.client, msg.Networks))

	case events.AccessPointAddedMsg:
		w.availableList.networks = append(w.availableList.networks, msg.AP)
		w.availableList.table.SetRows(AccessPointsToRows(w.availableList.networks))

	case events.AccessPointRemovedMsg:
		apPath := godbus.ObjectPath(msg.ApPath)
		result := slices.DeleteFunc(w.availableList.networks, func(n nm.AccessPoint) bool {
			return apPath == n.APPath
		})
		w.availableList.networks = result
		w.availableList.table.SetRows(AccessPointsToRows(result))

	case events.NewConnectionMsg:
		ap, _ := wifi.GetApFromSettings(
			w.client,
			godbus.ObjectPath(msg.ConnectionPath),
			w.availableList.networks,
		)
		w.savedList.networks = append(w.savedList.networks, *ap)
		w.savedList.table.SetRows(AccessPointsToRows(w.savedList.networks))

	case events.ConnectionRemovedMsg:
		conPath := godbus.ObjectPath(msg.ConnectionPath)
		result := slices.DeleteFunc(w.savedList.networks, func(n nm.AccessPoint) bool {
			return conPath == n.ConnectionPath
		})
		w.savedList.networks = result
		w.savedList.table.SetRows(AccessPointsToRows(result))

	case error:
		w.scanning = false
		w.err = msg
		return w, events.ShowError(msg)
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

func fetchSavedNetworks(c *dbus.Client, available []nm.AccessPoint) tea.Cmd {
	return func() tea.Msg {
		saved, err := wifi.GetSavedNetworks(c, available)
		if err != nil {
			return events.ShowError(err)
		}
		return events.SavedWifiMsg{Networks: saved}
	}
}

func fetchAvailableWifiNetworks(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := wifi.GetAvailableNetworks(c)
		if err != nil {
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}
