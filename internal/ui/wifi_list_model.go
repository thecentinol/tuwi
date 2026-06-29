package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/wifi"
	"strconv"
)

type WifiKeymap struct {
	connect,
	disconnect,
	forget,
	edit,
	autoConnect,
	scan key.Binding
}

type WifiListModel struct {
	client   *dbus.Client
	width    int
	height   int
	focused  bool
	networks []nm.AccessPoint
	table    comp.TableModel
	help     help.Model
	keymap   WifiKeymap
}

func NewSavedList(c *dbus.Client) WifiListModel {
	return WifiListModel{
		client: c,
		help:   help.New(),
		table: comp.NewTable(
			[]table.Column{
				// columns are set in root Model.sizeComponents()
				// so the widths can be calculated
			},
			[]table.Row{},
		),
		keymap: WifiKeymap{
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
				key.WithHelp("e", "edit-connection"),
			),
			autoConnect: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "auto-connect"),
			),
			scan: key.NewBinding(
				key.WithKeys(""),
				key.WithHelp("", ""),
				key.WithDisabled(),
			),
		},
	}
}

func NewAvailableList(c *dbus.Client) WifiListModel {
	return WifiListModel{
		client: c,
		help:   help.New(),
		table: comp.NewTable(
			[]table.Column{
				// columns ase set in root Model.sizeComponents()
				// so the widths can be calculated
			},
			[]table.Row{},
		),
		keymap: WifiKeymap{
			connect: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "connect"),
			),
			disconnect: key.NewBinding(
				key.WithKeys(""),
				key.WithHelp("", ""),
				key.WithDisabled(),
			),
			scan: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "scan"),
			),
			forget: key.NewBinding(
				key.WithKeys(""),
				key.WithHelp("", ""),
				key.WithDisabled(),
			),
			edit: key.NewBinding(
				key.WithKeys(""),
				key.WithHelp("", ""),
				key.WithDisabled(),
			),
			autoConnect: key.NewBinding(
				key.WithKeys(""),
				key.WithHelp("", ""),
				key.WithDisabled(),
			),
		},
	}
}

func (w WifiListModel) Update(msg tea.Msg) (WifiListModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, w.keymap.connect):
			w.ChooseConnection()

		case key.Matches(msg, w.keymap.disconnect):
			err := wifi.Disconnect(w.client)
			if err != nil {
				// TODO: call error modal.
			}

		case key.Matches(msg, w.keymap.scan):
			cmds = append(cmds, fetchAvailableWifiNetworks(w.client))
		}

	case events.WifiConnectReqMsg:
		err := wifi.ConnectToAvailableSecured(
			w.client,
			msg.Network,
			msg.Password,
		)

		if err != nil {
			return w, events.ShowError(err)
		}
	}

	w.table, cmd = w.table.Update(msg)
	cmds = append(cmds, cmd)

	return w, tea.Batch(cmds...)
}

func (w WifiListModel) View() tea.View {
	borderColor := base
	if w.focused {
		borderColor = focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(w.width).
		Height(w.height)

	return tea.NewView(style.Render(w.table.View().Content))
}

func (w WifiListModel) HelpView() []key.Binding {
	// append the wifi keybinds help view to the tables help.
	// table first because it has the up/down bindings.
	bindings := append(
		w.table.HelpView(),
		w.keymap.connect,
		w.keymap.disconnect,
		w.keymap.forget,
		w.keymap.edit,
		w.keymap.autoConnect,
		w.keymap.scan,
	)

	return bindings
}

func (w *WifiListModel) SelectedNetwork() *nm.AccessPoint {
	idx := w.table.Cursor()

	if idx < 0 || idx >= len(w.networks) {
		return nil
	}
	return &w.networks[idx]
}

func (w *WifiListModel) ChooseConnection() (WifiListModel, tea.Cmd) {
	selected := w.SelectedNetwork()
	switch {
	// Connect to saved network.
	case selected.IsSaved:
		wifi.ConnectToSaved(w.client, *selected)

	// Connect to available network that is not open.
	case selected.Secured && !selected.IsSaved:
		return *w, func() tea.Msg {
			return events.ShowPasswordModalMsg{Network: selected}
		}

	// Connect to available network that is open.
	case !selected.IsSaved && !selected.Secured && selected.SecurityType == "open":
		// TODO: connect to open.
	}

	return *w, nil
}

func AccessPointsToRows(networks []nm.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))
	for _, ap := range networks {
		rows = append(rows, table.Row{
			ap.SSID,
			ap.SecurityType,
			strconv.FormatBool(ap.Hidden),
			strconv.FormatInt(int64(ap.Strength), 10),
		})
	}
	return rows
}
