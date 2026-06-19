package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/models"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/wifi"
	"log"
)

type WifiKeymap struct {
	connect     key.Binding
	forget      key.Binding
	edit        key.Binding
	autoConnect key.Binding
	scan        key.Binding
}

type WifiListModel struct {
	client   *dbus.Client
	width    int
	height   int
	focused  bool
	networks []models.AccessPoint
	table    comp.TableModel
	help     help.Model
	keymap   WifiKeymap
}

type showPasswordModalMsg struct {
	network *models.AccessPoint
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
			forget: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "forget"),
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
			selected := w.SelectedNetwork()
			if selected == nil {
				break
			}

			if selected.Secured {
				return w, func() tea.Msg {
					return showPasswordModalMsg{network: selected}
				}
			}

		case key.Matches(msg, w.keymap.scan):
			cmds = append(cmds, fetchAvailableWifiNetworks(w.client))
		}

	case wifiConnectReqMsg:
		err := wifi.ConnectToAvailableSecured(
			w.client,
			msg.network,
			msg.password,
		)
		log.Printf("Connection req sent with password: %v", err)

		if err != nil {
			log.Printf("Error connecting to network: %v", err)
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
		w.keymap.forget,
		w.keymap.edit,
		w.keymap.autoConnect,
		w.keymap.scan,
	)

	return bindings
}

func (w *WifiListModel) SelectedNetwork() *models.AccessPoint {
	idx := w.table.Cursor()

	if idx < 0 || idx >= len(w.networks) {
		return nil
	}
	return &w.networks[idx]
}
