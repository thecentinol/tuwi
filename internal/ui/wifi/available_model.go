package wifi

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"fmt"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/keybindings"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/ui/theme"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

type AvailableModel struct {
	client *dbus.Client
	state  *wifidomain.State
	theme  *theme.Theme

	// displayedNetworks is used for indexing when slecting
	// a network to connect to because relying on the state of
	// State.Available will produce bugs if some APs are filtered out.
	displayedNetworks []nm.AccessPoint
	Table             comp.TableModel

	width  int
	height int

	keys      keybindings.Keybindings
	help      help.Model
	IsFocused bool
}

func NewWifiAvailableModel(
	c *dbus.Client,
	state *wifidomain.State,
	keys keybindings.Keybindings,
	theme *theme.Theme,
) AvailableModel {
	return AvailableModel{
		client: c,
		state:  state,
		theme:  theme,
		Table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
			keys,
			theme,
		),
		keys: keys,
	}
}

func (a AvailableModel) Init() tea.Cmd {
	return func() tea.Msg {
		networks, err := wifidomain.GetAvailableNetworks(a.client)
		if err != nil {
			return events.ShowError(err)
		}
		return events.AvailableWifiMsg{Networks: networks}
	}
}

func (a AvailableModel) Update(msg tea.Msg) (AvailableModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, a.keys.WifiConnectAvailable.ToBubbles()):
			selected := a.selectedAvailableNetwork()
			if selected == nil {
				return a, events.ShowError(fmt.Errorf("selectedAvailableNetwork: no network selected"))
			}
			_, cmd := a.handleConnect(*selected)
			cmds = append(cmds, cmd)

		case key.Matches(msg, a.keys.WifiScan.ToBubbles()):
			cmds = append(cmds, handleScan(a.client))
		}

	case events.AvailableWifiMsg:
		nearbySaved, err := wifidomain.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifidomain.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointAddedMsg:
		a.state.AddAvailable(msg.AP)
		nearbySaved, err := wifidomain.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifidomain.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointRemovedMsg:
		a.state.RemoveAvailable(msg.ApPath)
		nearbySaved, err := wifidomain.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifidomain.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.WifiConnectReqMsg:
		err := wifidomain.ConnectSecured(
			a.client,
			msg.Network,
			msg.Password,
		)

		if err != nil {
			return a, events.ShowError(err)
		}

	case events.NewConnectionMsg:
		nearbySaved, err := wifidomain.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifidomain.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.Table.SetRows(setAvailableRows(aps))

	case events.ConnectionRemovedMsg:
		nearbySaved, err := wifidomain.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := wifidomain.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)

		a.Table.SetRows(setAvailableRows(aps))
	}

	a.Table, cmd = a.Table.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a AvailableModel) View() tea.View {
	borderColor := a.theme.Border
	borderContent := a.theme.BorderContent
	if a.IsFocused {
		borderColor = a.theme.BorderFocused
		borderContent = a.theme.BorderContentFocused
	}

	title := comp.NewBorderContent(borderColor, borderContent, a.width, a.height)
	return tea.NewView(title.Render("[2]-Available networks", a.Table.View().Content, a.width))
}

func (a AvailableModel) HelpView() []key.Binding {
	bindings := append(
		a.Table.HelpView(),
		a.keys.WifiConnectAvailable.ToBubbles(),
	)
	return bindings
}

func (a *AvailableModel) Setheight(height int) {
	a.height = height
}

func (a *AvailableModel) SetWidth(width int) {
	a.width = width
	a.setAvailableColumns()
}

func (a *AvailableModel) selectedAvailableNetwork() *nm.AccessPoint {
	idx := a.Table.Cursor()

	if idx < 0 || idx >= len(a.state.Available) {
		return nil
	}
	return &a.displayedNetworks[idx]
}

func (a *AvailableModel) handleConnect(connection nm.AccessPoint) (AvailableModel, tea.Cmd) {
	securityType := wifidomain.DetermineSecurityType(connection.Flags, connection.WpaFlags, connection.RsnFlags)
	isSecured := securityType != "open" && securityType != "unknown"
	switch {
	case isSecured:
		return *a, func() tea.Msg {
			return events.ShowPasswordModalMsg{Network: &connection}
		}

	case securityType != "open" && securityType != "unknown":
		err := wifidomain.ConnectOpen(a.client, &connection)
		if err != nil {
			return *a, events.ShowError(err)
		}
	}

	return *a, nil
}

func (a *AvailableModel) setAvailableColumns() {
	colWidth := a.width / 5

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Freq", Width: colWidth},
		{Title: "Signal", Width: colWidth},
		{Title: "Channel", Width: colWidth - 10},
	}
	a.Table.SetColumns(cols)
}
