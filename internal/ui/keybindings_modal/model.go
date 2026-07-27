package keybindingsmodal

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/events"
	kb "github.com/thecentinol/tuwi/internal/keybindings"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type Model struct {
	theme *theme.Theme
	vp    viewport.Model

	keys kb.Keybindings
	help help.Model

	Height, Width int
	X, Y          int
	IsFocused     bool
}

func NewKeybindsModal(theme *theme.Theme, keys kb.Keybindings) Model {
	m := Model{
		theme: theme,
		vp:    viewport.New(),
		keys:  keys,
		help:  help.New(),
	}

	m.renderContent()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up.ToBubbles()):
			m.vp.ScrollUp(3)

		case key.Matches(msg, m.keys.Down.ToBubbles()):
			m.vp.ScrollDown(3)

		case key.Matches(msg, m.keys.GotoTop.ToBubbles()):
			m.vp.GotoTop()
			return m, nil

		case key.Matches(msg, m.keys.GotoBottom.ToBubbles()):
			m.vp.GotoBottom()
			return m, nil

		case key.Matches(msg, m.keys.Cancel.ToBubbles()):
			return m, events.CloseKeybindsModal()
		}
	}

	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	container := comp.NewBorderContent(
		m.theme.KeybindsModal.Border,
		m.theme.KeybindsModal.BorderContent,
		m.Height,
	)

	return tea.NewView(container.Render("Keybindings", m.vp.View(), m.Width))
}

func (m Model) HelpView() []key.Binding {
	// Note: OpenKeybindingsModal is excluded because it is handled
	// by the root model's HelpView.

	// Override default help descriptions to provide modal-specific context.
	cancel := m.keys.Cancel.ToBubbles()
	scrollUp := m.keys.Up.ToBubbles()
	scrollDown := m.keys.Down.ToBubbles()
	cancel.SetHelp(m.keys.Cancel.HelpKey, "Close")
	scrollUp.SetHelp(m.keys.Up.HelpKey, "Scroll up")
	scrollDown.SetHelp(m.keys.Down.HelpKey, "Scroll down")

	binds := []key.Binding{
		scrollUp,
		scrollDown,
		m.keys.GotoTop.ToBubbles(),
		m.keys.GotoBottom.ToBubbles(),
		cancel,
	}
	return binds
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.vp.SetWidth(width)
	m.vp.SetHeight(height)
}

func (m *Model) renderContent() {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.KeybindsModal.GroupText).
		MarginBottom(1)

	groupStyle := lipgloss.NewStyle().
		PaddingLeft(5).
		MarginBottom(2)

	groups := m.keys.AllGrouped()

	// styling the bindings
	help := m.help
	help.Styles.FullKey = lipgloss.NewStyle().
		Foreground(m.theme.KeybindsModal.KeybindsText)

	help.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(m.theme.KeybindsModal.DescriptionText)

	for _, group := range groups {
		title := titleStyle.Render("── " + group.Title + " ──")

		bindings := help.FullHelpView([][]key.Binding{
			kb.ToBubblesBatch(group.Bindings),
		})

		section := groupStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				bindings,
			),
		)

		// WriteString returns a nil error, the `revive` linter will complain about not handling the
		// error's for the below WriteString's. So just like the liberals, we tell em to shut up.
		builder.WriteString(section) //nolint:revive
		builder.WriteString("\n")    //nolint:revive
	}

	m.vp.SetContent(builder.String())
}
