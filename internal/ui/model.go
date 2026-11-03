package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/keybindings"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	kbmodal "github.com/thecentinol/tuwi/internal/ui/keybindings_modal"
	"github.com/thecentinol/tuwi/internal/ui/theme"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
)

type Focused string

const (
	FocusWifiSaved        Focused = "wifi-saved"
	FocusWifiAvailable    Focused = "wifi-available"
	FocusKeybindingsModal Focused = "keybindings-modal"
	FocusPasswordModal    Focused = "password-modal"
	FocusErrorModal       Focused = "error-modal"
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

	keybindingsModal     kbmodal.Model
	showKeybindingsModal bool

	passwordModal     comp.PasswordModel
	showPasswordModal bool

	errorModal     comp.ErrorModel
	showErrorModal bool

	startupErr error
}

func NewModel(client *dbus.Client, cfg config.AppConfig) Model {
	state := &wifidomain.State{}

	return Model{
		Client:           client,
		focus:            FocusWifiSaved,
		help:             help.New(),
		wifiSaved:        wifiui.NewWifiSavedModel(client, state, cfg),
		wifiAvailable:    wifiui.NewWifiAvailableModel(client, state, cfg),
		keybindingsModal: kbmodal.NewKeybindsModal(cfg),
		passwordModal:    comp.NewPasswordModal(cfg),
		errorModal:       comp.NewErrorModal(cfg),
		keys:             cfg.Keybinds,
		theme:            &cfg.Theme,
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

	case events.OpenKeybindsModalMsg:
		m.showKeybindingsModal = true
		return m, m.keybindingsModal.Init()

	case events.CloseKeybindsModalMsg:
		m.showKeybindingsModal = false
		// TODO: get previous focus before the keybindingsModal was opened and switch back to that.
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
		if m.showKeybindingsModal {
			m.keybindingsModal, cmd = m.keybindingsModal.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.wifiSaved.IsFocused = m.focus == FocusWifiSaved
	m.wifiAvailable.IsFocused = m.focus == FocusWifiAvailable
	m.keybindingsModal.IsFocused = m.focus == FocusKeybindingsModal

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var view tea.View
	wifiSavedView := m.wifiSaved.View().Content
	wifiAvailableView := m.wifiAvailable.View().Content
	kbView := m.keybindingsModal.View().Content
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

	if m.showKeybindingsModal && !m.showErrorModal && !m.showPasswordModal {
		kbModal := lipgloss.NewLayer(kbView).
			Z(1).
			X(m.keybindingsModal.X).
			Y(m.keybindingsModal.Y)

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

	case key.Matches(msg, m.keys.OpenKeybindingsModal.ToBubbles()):
		m.focus = FocusKeybindingsModal
		return events.OpenKeybindsModal()

	case m.showKeybindingsModal:
		m.keybindingsModal, cmd = m.keybindingsModal.Update(msg)

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
	if m.showKeybindingsModal {
		help = append(help, m.keybindingsModal.HelpView()...)
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

	if !m.showKeybindingsModal && !m.showErrorModal && !m.showPasswordModal {
		help = append(help, m.keys.OpenKeybindingsModal.ToBubbles())
	}

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
