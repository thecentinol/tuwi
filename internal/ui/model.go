package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	apphelp "github.com/thecentinol/tuwi/internal/ui/help"
	"github.com/thecentinol/tuwi/internal/ui/wifi"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
)

type Focused string

const (
	FocusWifiSaved     Focused = "wifi-saved"
	FocusWifiAvailable Focused = "wifi-available"
	FocusHelpMenu      Focused = "help-menu"
	FocusPasswordModal Focused = "password-modal"
	FocusErrorModal    Focused = "error-modal"
)

type Model struct {
	Client *dbus.Client
	width  int
	height int
	focus  Focused
	help   help.Model
	keys   config.Keybindings
	theme  *config.Theme

	wifiSaved       wifiui.SavedModel
	wifiAvailable   wifiui.AvailableModel
	selectedNetwork *nm.AccessPoint

	HelpMenu     apphelp.Model
	showHelpMenu bool

	passwordModal     comp.PasswordModel
	showPasswordModal bool

	errorModal     comp.ErrorModel
	showErrorModal bool

	startupErr error
}

func NewModel(client *dbus.Client, cfg config.AppConfig) Model {
	state := &wifi.State{}

	return Model{
		Client:        client,
		focus:         FocusWifiSaved,
		help:          help.New(),
		wifiSaved:     wifiui.NewWifiSavedModel(client, state, cfg),
		wifiAvailable: wifiui.NewWifiAvailableModel(client, state, cfg),
		HelpMenu:      apphelp.NewHelpMenu(cfg),
		passwordModal: comp.NewPasswordModal(cfg),
		errorModal:    comp.NewErrorModal(cfg),
		keys:          cfg.Keybinds,
		theme:         &cfg.Theme,
	}
}

func (m Model) SetStartupError(err error) Model {
	m.startupErr = err
	return m
}

func (m Model) showStartupError() tea.Cmd {
	if m.startupErr != nil {
		return events.ShowError(m.startupErr)
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.showStartupError(),
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
		m.setDimensions()

	case events.OpenHelpMenuMsg:
		m.showHelpMenu = true
		return m, m.HelpMenu.Init()

	case events.CloseHelpMenuMsg:
		m.showHelpMenu = false
		// TODO: get previous focus before the HelpMenu was opened and switch back to that.
		m.focus = FocusWifiAvailable

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
		m.showErrorModal = true
		m.focus = FocusErrorModal
		m.errorModal.SetText(msg.Err.Error())
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
		if m.showHelpMenu {
			m.HelpMenu, cmd = m.HelpMenu.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.wifiSaved.IsFocused = m.focus == FocusWifiSaved
	m.wifiAvailable.IsFocused = m.focus == FocusWifiAvailable
	m.HelpMenu.IsFocused = m.focus == FocusHelpMenu

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var view tea.View
	wifiSavedView := m.wifiSaved.View().Content
	wifiAvailableView := m.wifiAvailable.View().Content
	kbView := m.HelpMenu.View().Content
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

	if m.showHelpMenu && !m.showErrorModal && !m.showPasswordModal {
		kbModal := lipgloss.NewLayer(kbView).
			Z(1).
			X(m.HelpMenu.X).
			Y(m.HelpMenu.Y)

		layers = append(layers, kbModal)
	}

	comp := lipgloss.NewCompositor(layers...)
	view.SetContent(comp.Render())
	view.BackgroundColor = m.theme.BG
	view.AltScreen = true
	return view
}

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case m.showPasswordModal:
		m.passwordModal, cmd = m.passwordModal.Update(msg)

	case key.Matches(msg, m.keys.Quit.ToBubbles()):
		return tea.Quit

	case m.showErrorModal:
		m.focus = FocusErrorModal
		m.errorModal, cmd = m.errorModal.Update(msg)

	case key.Matches(msg, m.keys.OpenHelpMenu.ToBubbles()):
		m.focus = FocusHelpMenu
		return events.OpenHelpMenu()

	case m.showHelpMenu:
		m.HelpMenu, cmd = m.HelpMenu.Update(msg)

	case key.Matches(msg, m.keys.FocusedWifiSaved.ToBubbles()):
		m.focus = FocusWifiSaved

	case key.Matches(msg, m.keys.FocusedWifiAvailable.ToBubbles()):
		m.focus = FocusWifiAvailable

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
	if m.showHelpMenu {
		help = append(help, m.HelpMenu.HelpView()...)
	}
	if m.focus == FocusErrorModal {
		help = append(help, m.errorModal.HelpView()...)
	}
	if m.focus == FocusWifiSaved {
		help = append(help, m.wifiSaved.HelpView()...)
	}
	if m.focus == FocusWifiAvailable && !m.showPasswordModal {
		help = append(help, m.wifiAvailable.HelpView()...)
	}

	help = append(help, m.keys.Quit.ToBubbles())

	if !m.showHelpMenu && !m.showErrorModal && !m.showPasswordModal {
		help = append(help, m.keys.OpenHelpMenu.ToBubbles())
	}

	render := m.help
	render.Styles.Ellipsis = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.Ellipsis)

	render.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.ShortKey)

	render.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.ShortDesc)

	render.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.ShortSeparator)

	render.Styles.FullKey = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.FullKey)

	render.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.FullDesc)

	render.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(m.theme.HelpView.FullSeparator)

	return render.ShortHelpView(help)
}
