package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/wifi"
	"strings"
)

type WifiKeymap struct {
	up          key.Binding
	down        key.Binding
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
	networks []wifi.AccessPoint
	help     help.Model
	keymap   WifiKeymap
	cursor   int
}

func NewSavedList(c *dbus.Client) WifiListModel {
	return WifiListModel{
		client: c,
		help:   help.New(),
		keymap: WifiKeymap{
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
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
		keymap: WifiKeymap{
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
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

	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch {

		case key.Matches(msg, w.keymap.up):
			if w.cursor > 0 {
				w.cursor--
			}

		case key.Matches(msg, w.keymap.down):
			if w.cursor < len(w.networks)-1 {
				w.cursor++
			}

		case key.Matches(msg, w.keymap.scan):
			cmds = append(cmds, fetchAvailableWifiNetworks(w.client))
		}
	}

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

	var sb strings.Builder
	for i, v := range w.networks {
		if i == w.cursor && w.focused {
			sb.WriteString(focusedLine.Render(v.SSID))
		} else {
			sb.WriteString(v.SSID)
		}
		sb.WriteString("\n")
	}

	return tea.NewView(style.Render(sb.String()))
}

func (w WifiListModel) HelpView() string {
	return w.help.ShortHelpView([]key.Binding{
		w.keymap.up,
		w.keymap.down,
		w.keymap.connect,
		w.keymap.forget,
		w.keymap.edit,
		w.keymap.autoConnect,
		w.keymap.scan,
	})
}
