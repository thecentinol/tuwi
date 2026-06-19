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
	Input    textinput.Model
	password string
	keys     passwordKeymap

	X       int // x-coordinate for layering in root model
	Y       int // y-coordinate for layering in root model
	Width   int
	Focused bool
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
	// i.SetWidth(20)
	i.EchoMode = textinput.EchoPassword

	s := i.Styles()
	s.Cursor.Color = lipgloss.Color("205")
	s.Cursor.Blink = true
	// s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	i.SetStyles(s)

	return PasswordModel{
		Input: i,
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
			if p.Input.EchoMode == textinput.EchoPassword {
				p.Input.EchoMode = textinput.EchoNormal
			} else {
				p.Input.EchoMode = textinput.EchoPassword
			}
		case key.Matches(msg, p.keys.submit):
			password := p.Input.Value()
			p.Input.Reset()
			return p, handlePasswordSubmit(password)
		case key.Matches(msg, p.keys.cancel):
			p.Input.Reset()
			return p, handlePasswordClose()
		}
	}

	p.Input, cmd = p.Input.Update(msg)
	return p, cmd
}

func (p PasswordModel) View() tea.View {
	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Green).
		Width(p.Width)

	input := lipgloss.JoinVertical(
		lipgloss.Top,
		p.Input.View(),
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
