package info

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/config"
)

type InfoTableModel struct {
	Table table.Model
	theme *config.Theme
	keys  config.Keybindings
	help  help.Model
}

func NewInfoTable(
	cols []table.Column,
	rows []table.Row,
	keys config.Keybindings,
	theme *config.Theme,
) InfoTableModel {
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Background(theme.Table.SelectedBG.Color).
		Foreground(theme.Table.SelectedFG.Color)

	t.SetStyles(s)

	return InfoTableModel{
		Table: t,
		theme: theme,
		keys:  keys,
		help:  help.New(),
	}
}

func (t *InfoTableModel) SetColumns(cols []table.Column) {
	t.Table.SetColumns(cols)
}

func (t *InfoTableModel) SetRows(rows []table.Row) {
	t.Table.SetRows(rows)
	if len(rows) > 0 && t.Cursor() < 0 {
		t.Table.SetCursor(0)
	}
}

func (t *InfoTableModel) SetHeight(height int) {
	t.Table.SetHeight(height)
}

func (t *InfoTableModel) SetWidth(width int) {
	t.Table.SetWidth(width)
}

func (t *InfoTableModel) Update(msg tea.Msg) (InfoTableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, t.keys.Up.ToBubbles()):
			t.Table.MoveUp(1)

		case key.Matches(msg, t.keys.Down.ToBubbles()):
			t.Table.MoveDown(1)

		case key.Matches(msg, t.keys.GotoTop.ToBubbles()):
			t.Table.GotoTop()

		case key.Matches(msg, t.keys.GotoBottom.ToBubbles()):
			t.Table.GotoBottom()
		}
	}

	t.Table, cmd = t.Table.Update(msg)
	return *t, cmd
}

func (t *InfoTableModel) View() tea.View {
	// this removes the header
	_, view, _ := strings.Cut(t.Table.View(), "\n")
	return tea.NewView(view)
}

func (t InfoTableModel) HelpView() []key.Binding {
	return []key.Binding{
		t.keys.Up.ToBubbles(),
		t.keys.Down.ToBubbles(),
		t.keys.GotoTop.ToBubbles(),
		t.keys.GotoBottom.ToBubbles(),
	}
}

func (t *InfoTableModel) Cursor() int {
	return t.Table.Cursor()
}
