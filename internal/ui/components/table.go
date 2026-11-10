package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/events"
)

type TableModel struct {
	table table.Model
	theme *config.Theme
	keys  config.Keybindings
	help  help.Model
}

func NewTable(
	cols []table.Column,
	rows []table.Row,
	keys config.Keybindings,
	theme *config.Theme,
) TableModel {
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Background(theme.Table.HeaderBG.Color).
		Foreground(theme.Table.HeaderFG.Color)

	s.Selected = s.Selected.
		Background(theme.Table.SelectedBG.Color).
		Foreground(theme.Table.SelectedFG.Color)

	t.SetStyles(s)

	return TableModel{
		table: t,
		theme: theme,
		keys:  keys,
		help:  help.New(),
	}
}

func (t *TableModel) SetColumns(cols []table.Column) {
	t.table.SetColumns(cols)
}

func (t *TableModel) SetRows(rows []table.Row) {
	t.table.SetRows(rows)
	if len(rows) > 0 && t.Cursor() < 0 {
		t.table.SetCursor(0)
	}
}

func (t *TableModel) SetHeight(height int) {
	t.table.SetHeight(height)
}

func (t *TableModel) SetWidth(width int) {
	t.table.SetWidth(width)
}

func (t *TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, t.keys.Up.ToBubbles()):
			t.table.MoveUp(1)
			return *t, t.SendCursorPosition(t.Cursor())

		case key.Matches(msg, t.keys.Down.ToBubbles()):
			t.table.MoveDown(1)
			return *t, t.SendCursorPosition(t.Cursor())

		case key.Matches(msg, t.keys.GotoTop.ToBubbles()):
			t.table.GotoTop()
			return *t, t.SendCursorPosition(t.Cursor())

		case key.Matches(msg, t.keys.GotoBottom.ToBubbles()):
			t.table.GotoBottom()
			return *t, t.SendCursorPosition(t.Cursor())
		}
	}

	t.table, cmd = t.table.Update(msg)
	return *t, cmd
}

func (t *TableModel) View() tea.View {
	return tea.NewView(t.table.View())
}

func (t TableModel) HelpView() []key.Binding {
	return []key.Binding{
		t.keys.Up.ToBubbles(),
		t.keys.Down.ToBubbles(),
		t.keys.GotoTop.ToBubbles(),
		t.keys.GotoBottom.ToBubbles(),
	}
}

func (t *TableModel) Cursor() int {
	return t.table.Cursor()
}

func (t *TableModel) SendCursorPosition(pos int) tea.Cmd {
	return func() tea.Msg {
		return events.TableCursorMovedMsg{Idx: pos}
	}
}
