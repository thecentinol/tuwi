package info

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	comp "github.com/thecentinol/tuwi/internal/ui/components"
)

type Model struct {
	client *dbus.Client
	theme  *config.Theme
	keys   config.Keybindings

	Fields []Field
	Table  InfoTableModel

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
		theme:  &cfg.Theme,
		keys:   cfg.Keybinds,
		Table: NewInfoTable(
			[]table.Column{},
			[]table.Row{},
			cfg.Keybinds,
			&cfg.Theme,
		),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	borderColor := m.theme.Border
	borderContent := m.theme.BorderContent
	if m.Focused {
		borderColor = m.theme.BorderFocused
		borderContent = m.theme.BorderContentFocused
	}
	title := comp.NewBorderContent(borderColor.Color, borderContent.Color, m.height)

	return tea.NewView(title.Render("[3]-Info", m.Table.View().Content, m.width))
}

func (m *Model) SetDimensions(w, h int) {
	m.width = w
	m.setColumnsBro()

	m.height = h

	m.Table.SetWidth(w)
	m.Table.SetHeight(h - 2)
}

func (m *Model) setColumnsBro() {
	labelWidth := ((m.width * 1) / 3) - 3
	valueWidth := ((m.width * 2) / 3) - 3

	cols := []table.Column{
		{Title: "", Width: labelWidth},
		{Title: "", Width: valueWidth},
	}
	m.Table.SetColumns(cols)
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

	rows := make([]table.Row, 0, len(fields))
	for _, f := range fields {
		rows = append(rows, table.Row{
			f.Label,
			f.Value,
		})
	}

	m.Table.SetRows(rows)
}
