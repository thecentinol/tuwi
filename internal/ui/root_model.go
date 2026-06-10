package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
)

type keymap struct {
	focus1, focus2, up, down, scan, quit key.Binding
}

type Model struct {
	Client *dbus.Client
	width  int
	height int
	focus  int // which component is focused/active
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
			client: client,
		},
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
			scan: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "scan"),
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

		case key.Matches(msg, m.keymap.up), key.Matches(msg, m.keymap.down):

			switch m.focus {

			case 0:
				var cmd tea.Cmd
				m.wifi.savedList, cmd = m.wifi.savedList.Update(msg)
				cmds = append(cmds, cmd)

			case 1:
				var cmd tea.Cmd
				m.wifi.availableList, cmd = m.wifi.availableList.Update(msg)
				cmds = append(cmds, cmd)
			}

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
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.scan,
		m.keymap.quit,
	})

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
	m.wifi.availableList.width = halfWidth
	m.wifi.availableList.height = m.height - 2
}
