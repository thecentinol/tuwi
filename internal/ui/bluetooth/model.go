package bluetooth

import (
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	width   int  //nolint:unused
	height  int  //nolint:unused
	focused bool //nolint:unused
	cursor  int  //nolint:unused
}

func (m Model) Init() tea.Cmd {
	return nil
}

//nolint:unused
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) View() tea.View {
	return tea.NewView("hello bt world.")
}
