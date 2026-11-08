package info

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
)

const (
	// SetContent
	indent     = 2 // left margin for every row
	contentGap = 2 // space between the label column and the values
)

type Model struct {
	client *dbus.Client
	vp     viewport.Model
	theme  *config.Theme
	keys   config.Keybindings
	Fields []Field

	width,
	height int

	Focused bool
}

type Field struct {
	Label,
	Value string
}

func NewInfoModel(
	c *dbus.Client,
	cfg config.AppConfig,
) Model {
	return Model{
		client: c,
		vp:     viewport.New(),
		theme:  &cfg.Theme,
		keys:   cfg.Keybinds,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 1)
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
		}
	}

	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	borderColor := m.theme.Border
	borderContent := m.theme.BorderContent
	if m.Focused {
		borderColor = m.theme.BorderFocused
		borderContent = m.theme.BorderContentFocused
	}
	title := comp.NewBorderContent(borderColor.Color, borderContent.Color, m.height)

	return tea.NewView(title.Render("[3]-Info", m.vp.View(), m.width))
}

func (m *Model) SetDimensions(w, h int) {
	m.width, m.height = w, h
	m.vp.SetWidth(w)
	m.vp.SetHeight(h - 2)
}

func (m Model) HelpView() []key.Binding {
	// Override default help descriptions
	scrollUp := m.keys.Up.ToBubbles()
	scrollDown := m.keys.Down.ToBubbles()
	scrollUp.SetHelp(m.keys.Up.HelpKey, "Scroll up")
	scrollDown.SetHelp(m.keys.Down.HelpKey, "Scroll down")

	bindings := []key.Binding{
		scrollUp,
		scrollDown,
		m.keys.GotoTop.ToBubbles(),
		m.keys.GotoBottom.ToBubbles(),
	}
	return bindings
}

func (m *Model) SetContent(fields []Field) {
	m.Fields = fields
	var builder strings.Builder

	labelW := 0
	for _, f := range m.Fields {
		labelW = max(labelW, lipgloss.Width(f.Label))
	}

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Success).
		Width(labelW + contentGap)

	valueStyle := lipgloss.NewStyle().
		Width(max(m.vp.Width()-labelW-contentGap, 1))

	fieldStyle := lipgloss.NewStyle().
		PaddingLeft(1)

	for _, f := range m.Fields {
		label := labelStyle.Render(f.Label)
		value := valueStyle.Render(f.Value)

		field := fieldStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				label,
				value,
			),
		)

		builder.WriteString(field) //nolint:revive // the length and nil error are redundant
		builder.WriteString("\n")  //nolint:revive // the length and nil error are redundant
	}

	m.vp.SetContent(builder.String())
}
