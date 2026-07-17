package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type passwordKeymap struct {
	togglePassword key.Binding
	submit         key.Binding
	cancel         key.Binding
}

type PasswordModel struct {
	theme    *theme.Theme
	Input    textinput.Model
	password string
	keys     passwordKeymap
	help     help.Model

	Content string

	X       int // x-coordinate for layering in root model
	Y       int // y-coordinate for layering in root model
	Width   int
	Focused bool
}

func NewPasswordModal(theme *theme.Theme) PasswordModel {
	i := textinput.New()
	i.Placeholder = "Enter Password"
	i.SetVirtualCursor(true)
	i.Focus()
	i.EchoMode = textinput.EchoPassword

	s := i.Styles()
	s.Focused.Text = lipgloss.NewStyle().
		Foreground(theme.Input.TextFocused)

	s.Focused.Placeholder = lipgloss.NewStyle().
		Foreground(theme.Input.PlaceholderFocused)

	s.Blurred.Text = lipgloss.NewStyle().
		Foreground(theme.Input.TextBlurred)

	s.Blurred.Placeholder = lipgloss.NewStyle().
		Foreground(theme.Input.PlaceholderBlurred)

	s.Cursor.Color = theme.Input.Cursor
	s.Cursor.Blink = true
	i.SetStyles(s)

	return PasswordModel{
		theme: theme,
		Input: i,
		keys: passwordKeymap{
			togglePassword: key.NewBinding(
				key.WithKeys("ctrl+t"),
				key.WithHelp("ctrl+t", "Toggle Password"),
			),
			submit: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "connect"),
			),
			cancel: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "cancel"),
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
	container := NewBorderContent(p.theme.Password.Border, p.theme.Password.BorderContent, p.Width, 0)

	input := lipgloss.JoinHorizontal(
		lipgloss.Left,
		p.Input.View(),
		passwordIcon(p.Input.EchoMode == textinput.EchoPassword),
	)

	return tea.NewView(container.Render(p.Content, input, p.Width))
}

func (p PasswordModel) HelpView() []key.Binding {
	help := []key.Binding{
		p.keys.togglePassword,
		p.keys.submit,
		p.keys.cancel,
	}
	return help
}

func handlePasswordSubmit(password string) tea.Cmd {
	return func() tea.Msg {
		return events.PasswordResultMsg{
			Cancelled: false,
			Password:  password,
		}
	}
}

func handlePasswordClose() tea.Cmd {
	return func() tea.Msg {
		return events.PasswordResultMsg{
			Cancelled: true,
			Password:  "",
		}
	}
}

func passwordIcon(hidden bool) string {
	icon := lipgloss.NewStyle().
		Foreground(lipgloss.Green).
		Render("◉")

	if hidden {
		return "○"
	}
	return icon
}
