package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type passwordKeymap struct {
	togglePassword key.Binding
	submit         key.Binding
	cancel         key.Binding
}

type PasswordModel struct {
	height   int
	width    int
	input    textinput.Model
	password string
	keys     passwordKeymap
	focused  bool
}

type PasswordResultMsg struct {
	Cancelled bool
	Password  string
}

func NewPasswordModal() PasswordModel {
	i := textinput.New()
	i.Placeholder = "Enter Password"
	i.SetVirtualCursor(false)
	i.Focus()
	i.SetWidth(20)
	i.EchoMode = textinput.EchoPassword

	s := i.Styles()
	s.Cursor.Color = lipgloss.Color("205")
	s.Cursor.Blink = true
	// s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	i.SetStyles(s)

	return PasswordModel{
		input: i,
		keys: passwordKeymap{
			togglePassword: key.NewBinding(
				key.WithKeys("ctrl+t"),
			),
			submit: key.NewBinding(
				key.WithKeys("enter"),
			),
			cancel: key.NewBinding(
				key.WithKeys("esc"),
			),
		},
	}
}

func (p PasswordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p PasswordModel) Update(msg tea.Msg) (PasswordModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, p.keys.togglePassword):
			if p.input.EchoMode == textinput.EchoPassword {
				p.input.EchoMode = textinput.EchoNormal
			} else {
				p.input.EchoMode = textinput.EchoPassword
			}
		case key.Matches(msg, p.keys.submit):
			password := p.input.Value()
			p.input.Reset()
			return p, handlePasswordSubmit(password)
		case key.Matches(msg, p.keys.cancel):
			p.input.Reset()
			return p, handlePasswordClose()
		}
	}

	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p PasswordModel) View() tea.View {
	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Green).
		Width(20)

	input := lipgloss.JoinVertical(
		lipgloss.Top,
		p.input.View(),
		"escape = quit",
	)

	v := tea.NewView(container.Render(input))
	return v
}

func handlePasswordSubmit(password string) tea.Cmd {
	return func() tea.Msg {
		return PasswordResultMsg{
			Cancelled: false,
			Password:  password,
		}
	}
}

func handlePasswordClose() tea.Cmd {
	return func() tea.Msg {
		return PasswordResultMsg{
			Cancelled: true,
			Password:  "",
		}
	}
}
