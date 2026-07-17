package root

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

type keymap struct {
	focus1,
	focus2,
	quit key.Binding
}

type Focused string

const (
	FocusSaved         Focused = "saved"
	FocusAvailable     Focused = "available"
	FocusPasswordModal Focused = "password"
	FocusErrorModal    Focused = "error"
)

type Model struct {
	Client *dbus.Client
	width  int
	height int
	focus  Focused
	help   help.Model
	keymap keymap

	wifiSaved       wifiui.SavedModel
	wifiAvailable   wifiui.AvailableModel
	selectedNetwork *nm.AccessPoint

	passwordModal     comp.PasswordModel
	showPasswordModal bool

	errorModal     comp.ErrorModel
	showErrorModal bool
}

func NewModel(client *dbus.Client) Model {
	state := &wifidomain.State{}
	return Model{
		Client:        client,
		focus:         FocusSaved,
		help:          help.New(),
		wifiSaved:     wifiui.NewWifiSavedModel(client, state),
		wifiAvailable: wifiui.NewWifiAvailableModel(client, state),
		passwordModal: comp.NewPasswordModal(),
		errorModal:    comp.NewErrorModal(),
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
		m.wifiSaved.Init(),
		m.wifiAvailable.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case m.showPasswordModal:
			m.passwordModal, cmd = m.passwordModal.Update(msg)
			cmds = append(cmds, cmd)
		case m.showErrorModal:
			m.errorModal, cmd = m.errorModal.Update(msg)
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.focus1):
			m.focus = FocusSaved
		case key.Matches(msg, m.keymap.focus2):
			m.focus = FocusAvailable
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case m.wifiSaved.IsFocused:
			m.wifiSaved, cmd = m.wifiSaved.Update(msg)
			cmds = append(cmds, cmd)
		case m.wifiAvailable.IsFocused:
			m.wifiAvailable, cmd = m.wifiAvailable.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sizeComponents()

	case events.ShowPasswordModalMsg:
		m.showPasswordModal = true
		m.selectedNetwork = msg.Network
		m.passwordModal.Content = string(msg.Network.SSID)
		return m, m.passwordModal.Init()

	case events.PasswordResultMsg:
		m.showPasswordModal = false

		if msg.Cancelled {
			return m, nil
		}

		return m, func() tea.Msg {
			return events.WifiConnectReqMsg{
				Network:  m.selectedNetwork,
				Password: msg.Password,
			}
		}
	case events.ShowErrorMsg:
		m.errorModal.SetText(msg.Err.Error())
		m.showErrorModal = true
		return m, m.errorModal.Init()

	case events.DismissErrorMsg:
		m.showErrorModal = false

	default:
		m.wifiSaved, cmd = m.wifiSaved.Update(msg)
		cmds = append(cmds, cmd)

		m.wifiAvailable, cmd = m.wifiAvailable.Update(msg)
		cmds = append(cmds, cmd)

		if m.showPasswordModal {
			m.passwordModal, cmd = m.passwordModal.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.showErrorModal {
			m.errorModal, cmd = m.errorModal.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.wifiSaved.IsFocused = m.focus == FocusSaved
	m.wifiAvailable.IsFocused = m.focus == FocusAvailable

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var view tea.View
	wifiSavedView := m.wifiSaved.View().Content
	wifiAvailableView := m.wifiAvailable.View().Content
	helpView := m.HelpView()
	passwordView := m.passwordModal.View().Content
	errorModalView := m.errorModal.View().Content

	wifiView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		wifiSavedView,
		wifiAvailableView,
	)

	base := lipgloss.NewLayer(
		lipgloss.JoinVertical(
			lipgloss.Left,
			wifiView,
			helpView,
		),
	).Z(0)

	layers := []*lipgloss.Layer{base}

	if m.showPasswordModal && !m.showErrorModal {
		passwordModal := lipgloss.NewLayer(passwordView).
			Z(1).
			X(m.passwordModal.X).
			Y(m.passwordModal.Y)

		layers = append(layers, passwordModal)
	}

	if m.showErrorModal && !m.showPasswordModal {
		errorModal := lipgloss.NewLayer(errorModalView).
			Z(1).
			X(m.errorModal.X).
			Y(m.errorModal.Y)

		layers = append(layers, errorModal)
	}

	comp := lipgloss.NewCompositor(layers...)
	view.SetContent(comp.Render())
	view.AltScreen = true
	return view
}

func (m *Model) sizeComponents() {
	halfWidth := m.width / 2
	// halfHeight := m.height / 2
	helpHeight := 2
	modalWidth := m.width / 3
	tableHeight := m.height - helpHeight - 2

	// set wifi model's width and height
	m.wifiSaved.SetWidth(halfWidth)
	m.wifiSaved.Setheight(m.height - helpHeight)
	m.wifiSaved.Table.SetWidth(halfWidth)
	m.wifiSaved.Table.SetHeight(tableHeight)
	m.wifiAvailable.SetWidth(halfWidth)
	m.wifiAvailable.Setheight(m.height - helpHeight)
	m.wifiAvailable.Table.SetWidth(halfWidth)
	m.wifiAvailable.Table.SetHeight(tableHeight)

	// Password Modal
	iconWidth := 2
	m.passwordModal.Width = modalWidth
	m.passwordModal.Input.SetWidth(modalWidth - iconWidth - 4)

	// get height and width of password modal
	renderedPasswordModal := m.passwordModal.View().Content
	getPasswordModalHeight := lipgloss.Height(renderedPasswordModal)
	getPasswordModalWidth := lipgloss.Width(renderedPasswordModal)

	// calculate X and Y coordinates of password modal
	m.passwordModal.X = (m.width / 2) - (getPasswordModalWidth / 2)
	m.passwordModal.Y = (m.height / 2) - (getPasswordModalHeight / 2)

	// error modal
	m.errorModal.Width = m.width / 3
	m.errorModal.MaxHeight = m.height / 5
	m.errorModal.X = (m.width / 2) - (m.errorModal.Width / 2)
	m.errorModal.Y = (m.height / 2) - (m.errorModal.Height / 2)
}

func (m Model) HelpView() string {
	var help []key.Binding

	if m.showPasswordModal {
		help = append(help, m.passwordModal.HelpView()...)
	}
	if m.focus == FocusSaved {
		help = append(help, m.wifiSaved.HelpView()...)
	}
	if m.focus == FocusAvailable {
		help = append(help, m.wifiAvailable.HelpView()...)
	}

	help = append(help, m.keymap.quit)
	return m.help.ShortHelpView(help)
}
