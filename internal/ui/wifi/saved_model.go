package wifi

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/ui/theme"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

type savedKeymap struct {
	connect,
	disconnect,
	forget,
	edit,
	autoConnect key.Binding
}

type SavedModel struct {
	client *dbus.Client
	state  *wifidomain.State

	Table comp.TableModel

	width  int
	height int

	keys      savedKeymap
	help      help.Model
	IsFocused bool
}

func NewWifiSavedModel(c *dbus.Client, state *wifidomain.State) SavedModel {
	return SavedModel{
		client: c,
		state:  state,
		Table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
		),
		keys: savedKeymap{
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

func (s SavedModel) Init() tea.Cmd {
	an, err := wifidomain.GetActiveNetwork(s.client)
	if err != nil {
		return events.ShowError(err)
	}
	if an != nil {
		s.state.SetActiveConnection(*an)
	}

	return func() tea.Msg {
		saved, err := wifidomain.GetSavedNetworks(s.client)
		if err != nil {
			return events.ShowError(err)
		}
		return events.SavedWifiMsg{Networks: saved}
	}
}

func (s SavedModel) Update(msg tea.Msg) (SavedModel, tea.Cmd) {
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
			_, err := wifidomain.ConnectSaved(s.client, *selected)
			if err != nil {
				return s, events.ShowError(err)
			}

		case key.Matches(msg, s.keys.disconnect):
			err := wifidomain.Disconnect(s.client)
			if err != nil {
				return s, events.ShowError(err)
			}

		case key.Matches(msg, s.keys.forget):
			selected := s.selectedSavedNetwork()
			if selected == nil {
				return s, events.ShowError(fmt.Errorf("selectedSavedNetwork: no network selected"))
			}
			err := wifidomain.Forget(s.client, selected.Connection.ConnectionPath)
			if err != nil {
				return s, events.ShowError(err)
			}
		}

	case events.SavedWifiMsg:
		s.state.SetSaved(msg.Networks)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.AvailableWifiMsg:
		s.state.SetAvailable(msg.Networks)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.NewConnectionMsg:
		conn, err := wifidomain.BuildNewConnection(s.client, msg.ConnectionPath)
		if err != nil {
			return s, events.ShowError(err)
		}
		s.state.AddSaved(*conn)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.ConnectionRemovedMsg:
		s.state.RemoveSaved(msg.ConnectionPath)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.UpdateActiveConnectionMsg:
		an, err := wifidomain.GetActiveNetwork(s.client)
		if err != nil {
			return s, events.ShowError(err)
		}
		if an != nil {
			s.state.SetActiveConnection(*an)
			s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))
		}

	case events.ClearActiveConnectionMsg:
		s.state.ClearActiveConnection()
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))
	}

	s.Table, cmd = s.Table.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s SavedModel) View() tea.View {
	borderColor := theme.Base
	if s.IsFocused {
		borderColor = theme.Focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(s.width).
		Height(s.height)

	return tea.NewView(style.Render(s.Table.View().Content))
}

func (s SavedModel) HelpView() []key.Binding {
	bindings := append(
		s.Table.HelpView(),
		s.keys.connect,
		s.keys.disconnect,
		s.keys.forget,
	)
	return bindings
}

func (s *SavedModel) Setheight(height int) {
	s.height = height
}

func (s *SavedModel) SetWidth(width int) {
	s.width = width
	s.setSavedColumns()
}

func (s *SavedModel) selectedSavedNetwork() *nm.NearbyConnection {
	idx := s.Table.Cursor()

	if idx < 0 || idx >= len(s.state.Nearby) {
		return nil
	}
	return &s.state.Nearby[idx]
}

func (s *SavedModel) setSavedColumns() {
	colWidth := s.width / 5

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Status", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Strength", Width: colWidth - 10},
	}
	s.Table.SetColumns(cols)
}
