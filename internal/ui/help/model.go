package help

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/events"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
)

type Model struct {
	theme *config.Theme
	vp    viewport.Model

	keys config.Keybindings
	help help.Model

	height, width int
	x, y          int
	IsFocused     bool
}

func NewHelpMenu(cfg config.AppConfig) Model {
	m := Model{
		theme: &cfg.Theme,
		vp:    viewport.New(),
		keys:  cfg.Keybinds,
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
			return m, events.CloseHelpMenu()
		}
	}

	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	container := comp.NewBorderContent(
		m.theme.HelpMenu.Border.Color,
		m.theme.HelpMenu.BorderContent.Color,
		m.height,
	)

	return tea.NewView(container.Render("Help", m.vp.View(), m.width))
}

func (m Model) HelpView() []key.Binding {
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

func (m *Model) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.vp.SetWidth(width)
	m.vp.SetHeight(height)
}

func (m *Model) SetOverlayPosition(overlayX, overlayY int) {
	m.x, m.y = overlayX, overlayY
}

func (m *Model) renderContent() {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.HelpMenu.GroupText.Color).
		MarginBottom(1)

	groupStyle := lipgloss.NewStyle().
		PaddingLeft(5).
		MarginBottom(2)

	groups := m.keys.AllGrouped()

	// styling the bindings
	help := m.help
	help.Styles.FullKey = lipgloss.NewStyle().
		Foreground(m.theme.HelpMenu.KeybindsText.Color)

	help.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(m.theme.HelpMenu.DescriptionText.Color)

	for _, group := range groups {
		title := titleStyle.Render("── " + group.Title + " ──")

		bindings := help.FullHelpView([][]key.Binding{
			config.ToBubblesBatch(group.Bindings),
		})

		section := groupStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				bindings,
			),
		)

		// WriteString returns a nil error, the `revive` linter will complain about not handling the
		// error's for the below WriteString's. So just like the liberals, we tell em to shut up
		builder.WriteString(section) //nolint:revive
		builder.WriteString("\n")    //nolint:revive
	}

	m.vp.SetContent(builder.String())
}

func (m *Model) GetXposition() int {
	return m.x
}

func (m *Model) GetYposition() int {
	return m.y
}
