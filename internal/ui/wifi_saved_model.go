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

type wifiSavedKeymap struct {
	connect,
	disconnect,
	forget,
	edit,
	autoConnect key.Binding
}

type WifiSavedModel struct {
	client *dbus.Client
	state  *wifi.State

	table comp.TableModel

	width  int
	height int

	keys      wifiSavedKeymap
	help      help.Model
	isFocused bool
}

func NewWifiSavedModel(c *dbus.Client, state *wifi.State) WifiSavedModel {
	return WifiSavedModel{
		client: c,
		state:  state,
		table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
		),
		keys: wifiSavedKeymap{
			connect: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "connect"),
			),
			disconnect: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "disconnect"),
			),
			forget: key.NewBinding(
				key.WithKeys("f"),
				key.WithHelp("f", "forget"),
			),
			edit: key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit"),
				key.WithDisabled(),
			),
			autoConnect: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "auto-connect"),
				key.WithDisabled(),
			),
		},
	}
}

func (s WifiSavedModel) Init() tea.Cmd {
	return func() tea.Msg {
		saved, err := wifi.GetSavedNetworks(s.client)
		if err != nil {
			return events.ShowError(err)
		}
		return events.SavedWifiMsg{Networks: saved}
	}
}

func (s WifiSavedModel) Update(msg tea.Msg) (WifiSavedModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.connect):
			selected := s.selectedSavedNetwork()
			if selected == nil {
				return s, events.ShowError(fmt.Errorf("selectedSavedNetwork: no network selected"))
			}
			_, err := wifi.ConnectSaved(s.client, *selected)
			if err != nil {
				return s, events.ShowError(err)
			}

		case key.Matches(msg, s.keys.disconnect):
			err := wifi.Disconnect(s.client)
			if err != nil {
				return s, events.ShowError(err)
			}

		case key.Matches(msg, s.keys.forget):
			selected := s.selectedSavedNetwork()
			if selected == nil {
				return s, events.ShowError(fmt.Errorf("selectedSavedNetwork: no network selected"))
			}
			err := wifi.Forget(s.client, selected.Connection.ConnectionPath)
			if err != nil {
				return s, events.ShowError(err)
			}
		}

	case events.SavedWifiMsg:
		s.state.SetSaved(msg.Networks)
		s.table.SetRows(setSavedRows(s.state.Nearby))

	case events.AvailableWifiMsg:
		s.state.SetAvailable(msg.Networks)
		s.table.SetRows(setSavedRows(s.state.Nearby))

	case events.NewConnectionMsg:
		conn, err := wifi.BuildNewConnection(s.client, msg.ConnectionPath)
		if err != nil {
			return s, events.ShowError(err)
		}
		s.state.AddSaved(*conn)
		s.table.SetRows(setSavedRows(s.state.Nearby))

	case events.ConnectionRemovedMsg:
		s.state.RemoveSaved(msg.ConnectionPath)
		s.table.SetRows(setSavedRows(s.state.Nearby))
	}

	s.table, cmd = s.table.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s WifiSavedModel) View() tea.View {
	borderColor := base
	if s.isFocused {
		borderColor = focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(s.width).
		Height(s.height)

	return tea.NewView(style.Render(s.table.View().Content))
}

func (s WifiSavedModel) HelpView() []key.Binding {
	bindings := append(
		s.table.HelpView(),
		s.keys.connect,
		s.keys.disconnect,
		s.keys.forget,
	)
	return bindings
}

func (s *WifiSavedModel) Setheight(height int) {
	s.height = height
}

func (s *WifiSavedModel) SetWidth(width int) {
	s.width = width
	s.setSavedColumns()
}

func (s *WifiSavedModel) selectedSavedNetwork() *nm.NearbyConnection {
	idx := s.table.Cursor()

	if idx < 0 || idx >= len(s.state.Nearby) {
		return nil
	}
	return &s.state.Nearby[idx]
}

func (s *WifiSavedModel) setSavedColumns() {
	colWidth := s.width / 5

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Status", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Strength", Width: colWidth - 10},
	}
	s.table.SetColumns(cols)
}

func setSavedRows(nc []nm.NearbyConnection) []table.Row {
	rows := make([]table.Row, 0, len(nc))

	for _, n := range nc {
		isNearby := n.AP != nil
		var security string
		var strength string

		if isNearby {
			security = wifi.DetermineSecurityType(n.AP.Flags, n.AP.WpaFlags, n.AP.RsnFlags)
			strength = strconv.FormatInt(int64(n.AP.Strength), 10)
		} else {
			security = n.Connection.KeyMgmt
			strength = "0"
		}

		rows = append(rows, table.Row{
			n.Connection.SSID,
			"", // intentionally nil
			security,
			strconv.FormatBool(n.Connection.Hidden),
			strength + "%",
		})
	}

	return rows
}
