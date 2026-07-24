package components

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/keybindings"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type ErrorModel struct {
	theme *theme.Theme
	Text  string

	Width     int
	Height    int
	MaxHeight int
	X         int
	Y         int

	keys    keybindings.Keybindings
	Focused bool
}

func NewErrorModal(theme *theme.Theme, keys keybindings.Keybindings) ErrorModel {
	return ErrorModel{
		theme: theme,
		keys:  keys,
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
		case key.Matches(msg, e.keys.ErrorDismiss.ToBubbles()):
			return e, events.DismissError()
		}
	}

	return e, cmd
}

func (e ErrorModel) View() tea.View {
	// wrap the text.
	wrapped := lipgloss.NewStyle().
		Foreground(e.theme.Error.Text).
		Width(e.Width - 2).
		Render(e.Text)

	container := NewBorderContent(e.theme.Error.Border, e.theme.Error.BorderContent, e.Height-2)
	return tea.NewView(container.Render("Error:", wrapped, e.Width))
}

// This is used in root Model's Update().
func (e *ErrorModel) SetText(text string) {
	e.Text = text
	wrapped := lipgloss.NewStyle().
		Width(e.Width - 2).
		Render(e.Text)

	lines := strings.Count(wrapped, "\n") + 1
	e.Height = min(lines+2, e.MaxHeight)
}

func (e ErrorModel) HelpView() []key.Binding {
	help := []key.Binding{
		e.keys.ErrorDismiss.ToBubbles(),
	}
	return help
}
