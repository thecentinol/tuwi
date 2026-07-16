package bluetooth

import (
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	width   int
	height  int
	focused bool
	// devices []BluetoothDevice // TODO: implement this
	cursor int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) View() tea.View {
	return tea.NewView("hello bt world.")
}
