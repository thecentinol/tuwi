package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
)

type keymap struct {
	focus1, focus2, quit key.Binding
}

type Model struct {
	Client *dbus.Client
	width  int
	height int
	focus  int
	help   help.Model
	keymap keymap
	wifi   WifiModel

	// TODO: implement the following:
	// bluetooth BtModel
	// passwdModal   bool
}

func NewModel(client *dbus.Client) Model {
	return Model{
		Client: client,
		help:   help.New(),
		wifi: WifiModel{
			client:        client,
			savedList:     NewSavedList(client),
			availableList: NewAvailableList(client),
		},
		keymap: keymap{
			focus1: key.NewBinding(
				key.WithKeys("1"),
			),
			focus2: key.NewBinding(
				key.WithKeys("2"),
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
		m.wifi.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keymap.focus1):
			m.focus = 0
		case key.Matches(msg, m.keymap.focus2):
			m.focus = 1
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sizeComponents()
	}

	m.wifi.savedList.focused = m.focus == 0
	m.wifi.availableList.focused = m.focus == 1

	var cmd tea.Cmd
	m.wifi, cmd = m.wifi.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	help := m.wifi.HelpView()

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.wifi.View().Content,
			help,
		),
	)
	v.AltScreen = true
	return v
}

func (m *Model) sizeComponents() {
	halfWidth := m.width / 2
	// halfHeight := m.height / 2

	m.wifi.savedList.width = halfWidth
	m.wifi.savedList.height = m.height - 2
	m.wifi.savedList.table.SetWidth(m.wifi.savedList.width)
	m.wifi.savedList.table.SetHeight(m.wifi.savedList.height)

	m.wifi.availableList.width = halfWidth
	m.wifi.availableList.height = m.height - 2
	m.wifi.availableList.table.SetWidth(m.wifi.availableList.width)
	m.wifi.availableList.table.SetHeight(m.wifi.availableList.height)

	// set the width of the columns for the wifi tables
	ssidW := int(float64(halfWidth) * 0.5)
	securedW := int(float64(halfWidth) * 0.25)
	strengthW := int(float64(halfWidth) * 0.25)

	m.wifi.savedList.table.SetColumns([]table.Column{
		{Title: "SSID", Width: ssidW},
		{Title: "Secured", Width: securedW},
		{Title: "Strength", Width: strengthW - 10},
	})
	m.wifi.availableList.table.SetColumns([]table.Column{
		{Title: "SSID", Width: ssidW},
		{Title: "Secured", Width: securedW},
		{Title: "Strength", Width: strengthW - 10},
	})
}
