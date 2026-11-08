package wifi

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/service"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/ui/info"
)

type AvailableModel struct {
	client *dbus.Client
	state  *State
	theme  *config.Theme

	// displayedNetworks is used for indexing when slecting
	// a network to connect to because relying on the state of
	// State.Available will produce bugs if some APs are filtered out.
	displayedNetworks []nm.AccessPoint
	Table             comp.TableModel

	width     int
	height    int
	frameSize [2]int

	keys      config.Keybindings
	IsFocused bool
}

func NewWifiAvailableModel(
	c *dbus.Client,
	state *State,
	cfg config.AppConfig,
) AvailableModel {
	return AvailableModel{
		client: c,
		state:  state,
		theme:  &cfg.Theme,
		Table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
			cfg.Keybinds,
			&cfg.Theme,
		),
		keys: cfg.Keybinds,
	}
}

func (a AvailableModel) Init() tea.Cmd {
	return func() tea.Msg {
		networks, err := service.GetAvailableNetworks(a.client)
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
		nearbySaved, err := service.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := service.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointAddedMsg:
		a.state.AddAvailable(msg.AP)
		nearbySaved, err := service.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := service.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.AccessPointRemovedMsg:
		a.state.RemoveAvailable(msg.ApPath)
		nearbySaved, err := service.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := service.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.displayedNetworks = aps
		a.Table.SetRows(setAvailableRows(a.displayedNetworks))

	case events.WifiConnectReqMsg:
		err := service.ConnectSecured(
			a.client,
			msg.Network,
			msg.Password,
		)
		if err != nil {
			return a, events.ShowError(err)
		}

	case events.NewConnectionMsg:
		nearbySaved, err := service.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := service.DisplayAvailableAPs(
			nearbySaved,
			a.state.Available,
		)
		a.Table.SetRows(setAvailableRows(aps))

	case events.ConnectionRemovedMsg:
		nearbySaved, err := service.GetNearbySavedNetworks(a.client)
		if err != nil {
			return a, events.ShowError(err)
		}

		aps := service.DisplayAvailableAPs(
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

	title := comp.NewBorderContent(borderColor.Color, borderContent.Color, a.height)
	h, v := title.ContentStyle.GetFrameSize()
	a.frameSize = [2]int{h, v}
	return tea.NewView(title.Render("[2]-Available networks", a.Table.View().Content, a.width))
}

func (a AvailableModel) HelpView() []key.Binding {
	bindings := append(
		a.Table.HelpView(),
		a.keys.WifiConnectAvailable.ToBubbles(),
		a.keys.WifiScan.ToBubbles(),
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
	securityType := service.DetermineSecurityType(connection.Flags, connection.WpaFlags, connection.RsnFlags)
	isSecured := securityType != "open" && securityType != "unknown"
	switch {
	case isSecured:
		return *a, func() tea.Msg {
			return events.ShowPasswordModalMsg{Network: &connection}
		}

	case securityType != "open" && securityType != "unknown":
		err := service.ConnectOpen(a.client, &connection)
		if err != nil {
			return *a, events.ShowError(err)
		}
	}

	return *a, nil
}

func (a *AvailableModel) setAvailableColumns() {
	colWidth := (a.width / 5) - 3

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Freq", Width: colWidth},
		{Title: "Signal", Width: colWidth},
		{Title: "Channel", Width: colWidth},
	}
	a.Table.SetColumns(cols)
}

func (a *AvailableModel) InfoFields() []info.Field {
	i := a.Table.Cursor()
	if i < 0 || i >= len(a.displayedNetworks) {
		return nil
	}

	ap := a.displayedNetworks[i]

	return []info.Field{
		{Label: "SSID", Value: service.FormatSSID(ap.SSID)},
		{Label: "Frequency", Value: fmt.Sprintf("%s (Ch %v)",
			service.FormatFrequencyLong(ap.Frequency),
			service.DetermineChannel(ap.Frequency),
		)},
		{Label: "BSSID", Value: string(ap.BSSID)},
		{
			Label: "Security",
			Value: service.DetermineSecurityType(ap.Flags, ap.WpaFlags, ap.RsnFlags),
		},
		{Label: "Mode", Value: service.FormatAccessPointMode(ap.Mode)},
		{Label: "Max Bitrate", Value: service.FormatMaxBitrate(ap.MaxBitrate)},
		{Label: "Bandwidth", Value: service.FormatBandwidth(ap.Bandwidth)},
		{Label: "Strength", Value: service.FormatStrength(ap.Strength)},
		{
			Label: "Last Seen",
			Value: service.FormatLastSeen(service.ClockBootTimeToTime(ap.LastSeen)),
		},
	}
}
