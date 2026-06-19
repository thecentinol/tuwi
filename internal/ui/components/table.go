package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/thecentinol/tuwi/internal/models"
	"strconv"
)

type TableModel struct {
	table table.Model
	keys  keymap
	help  help.Model
}

type keymap struct {
	up     key.Binding
	down   key.Binding
	top    key.Binding
	bottom key.Binding
}

func NewTable(cols []table.Column, rows []table.Row) TableModel {
	return TableModel{
		table: table.New(
			table.WithColumns(cols),
			table.WithRows(rows),
		),
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

func AccessPointsToRows(networks []models.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))
	for _, ap := range networks {
		rows = append(rows, table.Row{
			ap.SSID,
			ap.SecurityType,
			strconv.FormatBool(ap.Hidden),
			strconv.FormatInt(int64(ap.Strength), 10),
		})
	}
	return rows
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
