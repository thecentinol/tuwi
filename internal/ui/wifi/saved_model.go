package wifi

import (
	"fmt"
	"strconv"

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

type SavedModel struct {
	client *dbus.Client
	state  *State
	theme  *config.Theme

	Table comp.TableModel

	width  int
	height int

	keys      *config.Keybindings
	IsFocused bool
}

func NewWifiSavedModel(
	c *dbus.Client,
	state *State,
	cfg config.AppConfig,
) SavedModel {
	return SavedModel{
		client: c,
		state:  state,
		theme:  &cfg.Theme,
		Table: comp.NewTable(
			[]table.Column{},
			[]table.Row{},
			cfg.Keybinds,
			&cfg.Theme,
		),
		keys: &cfg.Keybinds,
	}
}

func (s SavedModel) Init() tea.Cmd {
	an, err := service.GetActiveNetwork(s.client)
	if err != nil {
		return events.ShowError(err)
	}
	if an != nil {
		s.state.SetActiveConnection(*an)
	}

	return func() tea.Msg {
		saved, err := service.GetSavedNetworks(s.client)
		if err != nil {
			return events.ShowError(err)
		}
		return events.SavedWifiMsg{Networks: saved}
	}
}

func (s SavedModel) Update(msg tea.Msg) (SavedModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 1)
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.WifiConnectSaved.ToBubbles()):
			selected := s.selectedSavedNetwork()
			if selected == nil {
				return s, events.ShowError(fmt.Errorf("selectedSavedNetwork: no network selected"))
			}
			if selected.Connection.ConnectionPath != s.state.ActiveConnection.Connection {
				_, err := service.ConnectSaved(s.client, *selected)
				if err != nil {
					return s, events.ShowError(err)
				}
			}

		case key.Matches(msg, s.keys.WifiDisconnect.ToBubbles()):
			err := service.Disconnect(s.client)
			if err != nil {
				return s, events.ShowError(err)
			}

		case key.Matches(msg, s.keys.WifiForget.ToBubbles()):
			selected := s.selectedSavedNetwork()
			if selected == nil {
				return s, events.ShowError(fmt.Errorf("selectedSavedNetwork: no network selected"))
			}
			err := service.Forget(s.client, selected.Connection.ConnectionPath)
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
		conn, err := service.BuildNewConnection(s.client, msg.ConnectionPath)
		if err != nil {
			return s, events.ShowError(err)
		}
		s.state.AddSaved(*conn)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.ConnectionRemovedMsg:
		s.state.RemoveSaved(msg.ConnectionPath)
		s.Table.SetRows(setSavedRows(s.state.Nearby, &s.state.ActiveConnection))

	case events.UpdateActiveConnectionMsg:
		an, err := service.GetActiveNetwork(s.client)
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
	borderColor := s.theme.Border
	borderContent := s.theme.BorderContent
	if s.IsFocused {
		borderColor = s.theme.BorderFocused
		borderContent = s.theme.BorderContentFocused
	}

	title := comp.NewBorderContent(borderColor.Color, borderContent.Color, s.height)

	return tea.NewView(title.Render("[1]-Saved connections", s.Table.View().Content, s.width))
}

func (s SavedModel) HelpView() []key.Binding {
	bindings := append(
		s.Table.HelpView(),
		s.keys.WifiConnectSaved.ToBubbles(),
		s.keys.WifiDisconnect.ToBubbles(),
		s.keys.WifiForget.ToBubbles(),
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
	colWidth := (s.width / 5) - 3

	cols := []table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Status", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Signal", Width: colWidth},
	}
	s.Table.SetColumns(cols)
}

func (s *SavedModel) InfoFields() []info.Field {
	i := s.Table.Cursor()
	if i < 0 || i >= len(s.state.Saved) {
		return nil
	}

	sc := s.state.Saved[i]
	// TODO: replace magic num 16 with const
	fields := make([]info.Field, 0, 16)

	add := func(label, value string) {
		fields = append(fields, info.Field{Label: label, Value: value})
	}

	addIf := func(label, value string) {
		if value != "" {
			add(label, value)
		}
	}

	add("ConnectionPath", string(sc.ConnectionPath))

	add("UUID", string(sc.UUID))
	add("Auto Connect", strconv.FormatBool(sc.AutoConnect))

	add("SSID", sc.SSID)
	add("Seen BSSIDs", service.FormatList(sc.BSSIDs))
	add("Mode", string(sc.Mode))
	add("Band", service.FormatSavedConnectionBand(sc.Band))
	add("Hidden", strconv.FormatBool(sc.Hidden))
	addIf("Mac Address", sc.MacAddress)

	add("Key Mgmt", sc.KeyMgmt)
	add("Auth Alg", string(sc.AuthAlg))
	add("Proto", service.FormatSavedConnectionProto(sc.Proto))
	add("PSK Flags", service.FormatSecretFlags(nm.SecretFlags(sc.PskFlags)))
	add("WEP Key Flags", service.FormatSecretFlags(nm.SecretFlags(sc.WepKeyFlags)))
	add("WEP Key Type", sc.WepKeyType.String())
	add("WPS Method", sc.WpsMethod.String())

	return fields
}
