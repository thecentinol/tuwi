package root

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/keybindings"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/ui/theme"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

type Focused string

// make sure these match Keybindings ViewName's
const (
	FocusSaved         Focused = "wifi-saved"
	FocusAvailable     Focused = "wifi-available"
	FocusPasswordModal Focused = "password-modal"
	FocusErrorModal    Focused = "error-modal"
)

type Model struct {
	Client *dbus.Client
	width  int
	height int
	focus  Focused
	help   help.Model
	keys   keybindings.Keybindings
	theme  *theme.Theme

	wifiSaved       wifiui.SavedModel
	wifiAvailable   wifiui.AvailableModel
	selectedNetwork *nm.AccessPoint

	passwordModal     comp.PasswordModel
	showPasswordModal bool

	errorModal     comp.ErrorModel
	showErrorModal bool
}

func NewModel(client *dbus.Client, keys keybindings.Keybindings) Model {
	th := theme.Default
	state := &wifidomain.State{}
	return Model{
		Client:        client,
		focus:         FocusSaved,
		help:          help.New(),
		wifiSaved:     wifiui.NewWifiSavedModel(client, state, keys, &th),
		wifiAvailable: wifiui.NewWifiAvailableModel(client, state, keys, &th),
		passwordModal: comp.NewPasswordModal(&th, keys),
		errorModal:    comp.NewErrorModal(&th, keys),
		keys:          keys,
		theme:         &th,
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
		cmd = m.handleKeyPressMsg(msg)
		cmds = append(cmds, cmd)

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
	view.BackgroundColor = m.theme.BG
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

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case m.showPasswordModal:
		m.passwordModal, cmd = m.passwordModal.Update(msg)

	case key.Matches(msg, m.keys.Quit.ToBubbles()):
		return tea.Quit

	case m.showErrorModal:
		m.errorModal, cmd = m.errorModal.Update(msg)

	case key.Matches(msg, m.keys.FocusedWifiSaved.ToBubbles()):
		m.focus = FocusSaved

	case key.Matches(msg, m.keys.FocusedWifiAvailable.ToBubbles()):
		m.focus = FocusAvailable

	case m.wifiSaved.IsFocused:
		m.wifiSaved, cmd = m.wifiSaved.Update(msg)

	case m.wifiAvailable.IsFocused:
		m.wifiAvailable, cmd = m.wifiAvailable.Update(msg)
	}

	return cmd
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

	// help = append(help, m.keymap.quit)
	render := m.help
	render.Styles.Ellipsis = lipgloss.NewStyle().
		Foreground(m.theme.Help.Ellipsis)

	render.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(m.theme.Help.ShortKey)

	render.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(m.theme.Help.ShortDesc)

	render.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(m.theme.Help.ShortSeparator)

	render.Styles.FullKey = lipgloss.NewStyle().
		Foreground(m.theme.Help.FullKey)

	render.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(m.theme.Help.FullDesc)

	render.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(m.theme.Help.FullSeparator)

	return render.ShortHelpView(help)
}
