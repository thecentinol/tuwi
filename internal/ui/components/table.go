package components

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/thecentinol/tuwi/internal/wifi"
)

type TableModel struct {
	table table.Model
}

func NewTable(cols []table.Column, rows []table.Row) TableModel {
	return TableModel{
		table: table.New(
			table.WithColumns(cols),
			table.WithRows(rows),
		),
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
	t.table, cmd = t.table.Update(msg)
	return *t, cmd
}

func (t *TableModel) View() tea.View {
	return tea.NewView(t.table.View())
}

func AccessPointsToRows(networks []wifi.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))
	for _, ap := range networks {
		rows = append(rows, table.Row{
			ap.SSID,
			fmt.Sprintf("%v", ap.Secured),
			fmt.Sprintf("%d\n", ap.Strength),
		})
	}
	return rows
}
