package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"strconv"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/wifi"
)

type wifiAvailableKeymap struct {
	connect,
	scan key.Binding
}

type WifiAvailableModel struct {
	client *dbus.Client
	state  *wifi.State

	// displayedNetworks is used for indexing when slecting
	// a network to connect to because relying on the state of
	// State.Available will produce bugs if some APs are filtered out.
	displayedNetworks []nm.AccessPoint
	table             comp.TableModel

	width  int
	height int

	keys      wifiAvailableKeymap
	help      help.Model
	isFocused bool
}

func NewWifiAvailableModel(c *dbus.Client, state *wifi.State) WifiAvailableModel {
	return WifiAvailableModel{
		client: c,
		state:  state,
		table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
		),
		keys: wifiAvailableKeymap{
			connect: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "connect"),
			),
			scan: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "scan"),
			),
		},
	}
}

func (a WifiAvailableModel) Init() tea.Cmd {
	return func() tea.Msg {
		networks, err := wifi.GetAvailableNetworks(a.client)
		if err != nil {
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}

func (a WifiAvailableModel) Update(msg tea.Msg) (WifiAvailableModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, a.keys.connect):
			selected := a.selectedAvailableNetwork()
			if selected == nil {
				return a, events.ShowError(fmt.Errorf("selectedAvailableNetwork: no network selected"))
			}
			_, cmd := a.handleConnect(*selected)
			cmds = append(cmds, cmd)

		case key.Matches(msg, a.keys.scan):
			cmds = append(cmds, handleScan(a.client))
		}

	case events.AvailableWifiMsg:
		nearbySaved, err := wifi.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifi.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointAddedMsg:
		a.state.AddAvailable(msg.AP)
		nearbySaved, err := wifi.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifi.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointRemovedMsg:
		a.state.RemoveAvailable(msg.ApPath)
		nearbySaved, err := wifi.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifi.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.WifiConnectReqMsg:
		err := wifi.ConnectSecured(
			a.client,
			msg.Network,
			msg.Password,
		)

		if err != nil {
			return a, events.ShowError(err)
		}

	case events.NewConnectionMsg:
		nearbySaved, err := wifi.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifi.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.table.SetRows(setAvailableRows(aps))

	case events.ConnectionRemovedMsg:
		nearbySaved, err := wifi.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifi.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)

		a.table.SetRows(setAvailableRows(aps))
	}

	a.table, cmd = a.table.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a WifiAvailableModel) View() tea.View {
	borderColor := base
	if a.isFocused {
		borderColor = focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(a.width).
		Height(a.height)

	return tea.NewView(style.Render(a.table.View().Content))
}

func (a WifiAvailableModel) HelpView() []key.Binding {
	bindings := append(
		a.table.HelpView(),
		a.keys.connect,
		a.keys.scan,
	)
	return bindings
}

func (a *WifiAvailableModel) Setheight(height int) {
	a.height = height
}

func (a *WifiAvailableModel) SetWidth(width int) {
	a.width = width
	a.setAvailableColumns()
}

func (a *WifiAvailableModel) selectedAvailableNetwork() *nm.AccessPoint {
	idx := a.table.Cursor()

	if idx < 0 || idx >= len(a.state.Available) {
		return nil
	}
	return &a.displayedNetworks[idx]
}

func (a *WifiAvailableModel) handleConnect(connection nm.AccessPoint) (WifiAvailableModel, tea.Cmd) {
	securityType := wifi.DetermineSecurityType(connection.Flags, connection.WpaFlags, connection.RsnFlags)
	isSecured := securityType != "open" && securityType != "unknown"
	switch {
	case isSecured:
		return *a, func() tea.Msg {
			return events.ShowPasswordModalMsg{Network: &connection}
		}

	case securityType != "open" && securityType != "unknown":
		err := wifi.ConnectOpen(a.client, &connection)
		if err != nil {
			return *a, events.ShowError(err)
		}
	}

	return *a, nil
}

func handleScan(c *dbus.Client) tea.Cmd {
	return func() tea.Msg {
		networks, err := wifi.GetAvailableNetworks(c)
		if err != nil {
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}

func (a *WifiAvailableModel) setAvailableColumns() {
	colWidth := a.width / 4

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Strength", Width: colWidth - 10},
	}
	a.table.SetColumns(cols)
}

func setAvailableRows(networks []nm.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))

	for _, n := range networks {
		strength := strconv.FormatInt(int64(n.Strength), 10)
		securityType := wifi.DetermineSecurityType(n.Flags, n.WpaFlags, n.RsnFlags)

		rows = append(rows, table.Row{
			wifi.FormatSSID(n.SSID),
			securityType,
			"",
			string(strength) + "%",
		})
	}

	return rows
}
