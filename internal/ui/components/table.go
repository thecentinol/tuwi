package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type keymap struct {
	up     key.Binding
	down   key.Binding
	top    key.Binding
	bottom key.Binding
}

type TableModel struct {
	table table.Model
	theme *theme.Theme
	keys  keymap
	help  help.Model
}

func NewTable(cols []table.Column, rows []table.Row, theme *theme.Theme) TableModel {
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Background(theme.Table.HeaderBG).
		Foreground(theme.Table.HeaderFG)

	s.Selected = s.Selected.
		Background(theme.Table.SelectedBG).
		Foreground(theme.Table.SelectedFG)

	t.SetStyles(s)

	return TableModel{
		table: t,
		theme: theme,
		keys: keymap{
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
			top: key.NewBinding(
				key.WithKeys("t"),
				key.WithHelp("t", "go-to-top"),
			),
			bottom: key.NewBinding(
				key.WithKeys("b"),
				key.WithHelp("b", "go-to-bottom"),
			),
		},
		help: help.New(),
	}
}

func (t *TableModel) SetColumns(cols []table.Column) {
	t.table.SetColumns(cols)
}

func (t *TableModel) SetRows(rows []table.Row) {
	t.table.SetRows(rows)
}

// the height and width gets set in root Model.sizeComponents
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
		case key.Matches(msg, t.keys.up):
			t.table.MoveUp(1)
		case key.Matches(msg, t.keys.down):
			t.table.MoveDown(1)
		case key.Matches(msg, t.keys.top):
			t.table.GotoTop()
		case key.Matches(msg, t.keys.bottom):
			t.table.GotoBottom()
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
		t.keys.up,
		t.keys.down,
		t.keys.top,
		t.keys.bottom,
	}
}

func (t *TableModel) Cursor() int {
	return t.table.Cursor()
}
