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
)

type keymap struct {
	focus1,
	focus2,
	quit key.Binding
}

type Model struct {
	Client            *dbus.Client
	width             int
	height            int
	focus             int
	help              help.Model
	keymap            keymap
	wifi              WifiModel
	passwordModal     comp.PasswordModel
	showPasswordModal bool
	selectedNetwork   *models.AccessPoint

	// TODO: implement the following:
	// bluetooth BtModel
	// passwdModal   bool
}

type wifiConnectReqMsg struct {
	network  *models.AccessPoint
	password string
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
		passwordModal: comp.NewPasswordModal(),
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
	var cmd tea.Cmd

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

	case showPasswordModalMsg:
		m.showPasswordModal = true
		m.selectedNetwork = msg.network
		return m, m.passwordModal.Init()

	case comp.PasswordResultMsg:
		m.showPasswordModal = false

		if msg.Cancelled {
			return m, nil
		}

		return m, func() tea.Msg {
			return wifiConnectReqMsg{
				network:  m.selectedNetwork,
				password: msg.Password,
			}
		}
	}

	if m.showPasswordModal {
		m.passwordModal, cmd = m.passwordModal.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.wifi.savedList.focused = m.focus == 0
		m.wifi.availableList.focused = m.focus == 1

		m.wifi, cmd = m.wifi.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var view tea.View
	wifiView := m.wifi.View().Content
	helpView := m.HelpView()
	passwordView := m.passwordModal.View().Content

	base := lipgloss.NewLayer(
		lipgloss.JoinVertical(
			lipgloss.Left,
			wifiView,
			helpView,
		),
	).Z(0)

	layers := []*lipgloss.Layer{base}

	if m.showPasswordModal {
		passwordModal := lipgloss.NewLayer(passwordView).
			Z(1).
			X(m.passwordModal.X).
			Y(m.passwordModal.Y)

		layers = append(layers, passwordModal)
	}

	comp := lipgloss.NewCompositor(layers...)
	view.SetContent(comp.Render())
	view.AltScreen = true
	return view
}

func (m *Model) sizeComponents() {
	halfWidth := m.width / 2
	// halfHeight := m.height / 2
	helpHeight := 1
	tableH := m.height - helpHeight - 2

	m.wifi.savedList.width = halfWidth
	m.wifi.savedList.height = m.height - helpHeight
	m.wifi.savedList.table.SetWidth(halfWidth)
	m.wifi.savedList.table.SetHeight(tableH)

	m.wifi.availableList.width = halfWidth
	m.wifi.availableList.height = m.height - helpHeight
	m.wifi.availableList.table.SetWidth(halfWidth)
	m.wifi.availableList.table.SetHeight(tableH)

	// set the width of the columns for the wifi tables
	colWidth := int(float64(halfWidth) * 0.25)

	m.wifi.savedList.table.SetColumns([]table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Strength", Width: colWidth - 10},
	})
	m.wifi.availableList.table.SetColumns([]table.Column{
		{Title: "SSID", Width: colWidth},
		{Title: "Security", Width: colWidth},
		{Title: "Hidden", Width: colWidth},
		{Title: "Strength", Width: colWidth - 10},
	})

	// Password Modal
	modalWidth := m.width / 3
	iconWidth := 2
	m.passwordModal.Width = modalWidth
	m.passwordModal.Input.SetWidth(modalWidth - iconWidth - 4)

	// get height and width of the modal
	renderedModal := m.passwordModal.View().Content
	getModalHeight := lipgloss.Height(renderedModal)
	getModalWidth := lipgloss.Width(renderedModal)

	// calculate X and Y coordinates
	m.passwordModal.X = (m.width / 2) - (getModalWidth / 2)
	m.passwordModal.Y = (m.height / 2) - (getModalHeight / 2)
}

func (m Model) HelpView() string {
	var help []key.Binding

	if m.showPasswordModal {
		help = append(help, m.passwordModal.HelpView()...)
	} else {
		help = append(help, m.wifi.HelpView()...)
	}

	help = append(help, m.keymap.quit)
	return m.help.ShortHelpView(help)
}
