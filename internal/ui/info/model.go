package info

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/service"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
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
	var cmds []tea.Cmd
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
