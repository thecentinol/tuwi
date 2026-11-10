package info

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
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

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.CopyInfoField.ToBubbles()):
			err := m.Copy()
			if err != nil {
				return m, events.ShowError(fmt.Errorf("failed to copy to clipboard: %w", err))
			}
		}
	}

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
	bindings := []key.Binding{
		m.keys.Up.ToBubbles(),
		m.keys.Down.ToBubbles(),
		m.keys.GotoTop.ToBubbles(),
		m.keys.GotoBottom.ToBubbles(),
		m.keys.CopyInfoField.ToBubbles(),
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

// copy the Value of a Field
func (m *Model) Copy() error {
	if len(m.Fields) == 0 {
		return nil
	}

	field := m.Fields[m.Table.Cursor()]
	if err := clipboard.WriteAll(field.Value); err != nil {
		return err
	}

	return nil
}
