package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
)

// this will render both known and unknown wifi-networks/bt-devices
type Model struct {
	Client                *dbus.Client
	width                 int
	height                int
	focus                 int // which component is focused/active
	savedWifiNetworks     WifiListModel
	availableWifiNetworks WifiListModel

	// TODO: implement the following:
	// knownBt       BtListModel
	// availableBt   BtListModel
	// passwdModal   bool
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.focus = 1
		case "2":
			m.focus = 2
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sizeComponents()
	}

	m.savedWifiNetworks.focused = m.focus == 1
	m.availableWifiNetworks.focused = m.focus == 2

	return m, nil
}

func (m Model) View() tea.View {
	v := tea.NewView(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.savedWifiNetworks.View().Content,
			m.availableWifiNetworks.View().Content,
		),
	)
	v.AltScreen = true
	return v
}

func (m *Model) sizeComponents() {
	halfWidth := m.width / 2
	// halfHeight := m.height / 2

	m.savedWifiNetworks.width = halfWidth
	m.savedWifiNetworks.height = m.height
	m.availableWifiNetworks.width = halfWidth
	m.availableWifiNetworks.height = m.height
}
