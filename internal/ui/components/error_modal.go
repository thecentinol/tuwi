package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"strings"

	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type ErrorModalKeymap struct {
	dismiss key.Binding
}

type ErrorModel struct {
	theme *theme.Theme
	Text  string

	Width     int
	Height    int
	MaxHeight int
	X         int
	Y         int

	keys    ErrorModalKeymap
	help    help.Model
	Focused bool
}

func NewErrorModal(theme *theme.Theme) ErrorModel {
	return ErrorModel{
		theme: theme,
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
		Foreground(e.theme.Error.Text).
		Width(e.Width - 2).
		Render(e.Text)

	conatiner := NewBorderContent(lipgloss.Green, e.Width, e.Height-2)
	return tea.NewView(conatiner.Render("Error:", wrapped, e.Width))
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
		e.keys.dismiss,
	}
	return help
}
