package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
)

type keymap struct {
	focus1, focus2, up, down, quit key.Binding
}

// this will render both known and unknown wifi-networks/bt-devices
type Model struct {
	Client                *dbus.Client
	width                 int
	height                int
	focus                 int // which component is focused/active
	help                  help.Model
	keymap                keymap
	savedWifiNetworks     WifiListModel
	availableWifiNetworks WifiListModel

	// TODO: implement the following:
	// knownBt       BtListModel
	// availableBt   BtListModel
	// passwdModal   bool
}

func NewModel(client *dbus.Client) Model {
	return Model{
		Client: client,
		help:   help.New(),
		keymap: keymap{
			focus1: key.NewBinding(
				key.WithKeys("1"),
			),
			focus2: key.NewBinding(
				key.WithKeys("2"),
			),
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchSavedNetworks(m.Client),
		fetchAvailableWifiNetworks(m.Client),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case savedWifiMsg:
		m.savedWifiNetworks.networks = msg
		return m, nil
	case availableWifiMsg:
		m.availableWifiNetworks.networks = msg
		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keymap.focus1):
			m.focus = 0
		case key.Matches(msg, m.keymap.focus2):
			m.focus = 1
		case key.Matches(msg, m.keymap.up), key.Matches(msg, m.keymap.down):
			switch m.focus {
			case 0:
				m.savedWifiNetworks, _ = m.savedWifiNetworks.Update(msg)
			case 1:
				m.availableWifiNetworks, _ = m.availableWifiNetworks.Update(msg)
			}
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sizeComponents()
	}

	m.savedWifiNetworks.focused = m.focus == 0
	m.availableWifiNetworks.focused = m.focus == 1

	return m, nil
}

func (m Model) View() tea.View {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.quit,
	})

	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.savedWifiNetworks.View().Content,
		m.availableWifiNetworks.View().Content,
	)

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			top,
			help,
		),
	)
	v.AltScreen = true
	return v
}

func (m *Model) sizeComponents() {
	halfWidth := m.width / 2
	// halfHeight := m.height / 2

	m.savedWifiNetworks.width = halfWidth
	m.savedWifiNetworks.height = m.height - 2
	m.availableWifiNetworks.width = halfWidth
	m.availableWifiNetworks.height = m.height - 2
}
