package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thecentinol/tuwi/internal/events"
	"strings"
)

type ErrorModalKeymap struct {
	dismiss key.Binding
}

type ErrorModel struct {
	Text string

	Width     int
	Height    int
	MaxHeight int
	X         int
	Y         int

	keys    ErrorModalKeymap
	help    help.Model
	Focused bool
}

func NewErrorModal() ErrorModel {
	return ErrorModel{
		keys: ErrorModalKeymap{
			dismiss: key.NewBinding(
				key.WithKeys("enter", "esc"),
				key.WithHelp("enter/esc", "dismiss"),
			),
		},
	}
}

func (e ErrorModel) Init() tea.Cmd {
	return nil
}

func (e ErrorModel) Update(msg tea.Msg) (ErrorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, e.keys.dismiss):
			return e, events.DismissError()
		}
	}

	return e, cmd
}

func (e ErrorModel) View() tea.View {
	// wrap the text.
	wrapped := lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		Width(e.Width - 2).
		Render(e.Text)

	// This is the border around the text content.
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Green).
		Width(e.Width).
		Height(e.Height - 2)

	view := tea.NewView(style.Render(wrapped))
	return view
}

// This is used in root Model's Update().
func (e *ErrorModel) SetText(text string) {
	e.Text = text
	wrapped := lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		Width(e.Width - 2).
		Render(e.Text)

	lines := strings.Count(wrapped, "\n") + 1
	e.Height = min(lines+2, e.MaxHeight)
}

func (e ErrorModel) HelpView() []key.Binding {
	help := []key.Binding{
		e.keys.dismiss,
	}
	return help
}
