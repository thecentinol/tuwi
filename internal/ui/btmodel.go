package ui

import (
	tea "charm.land/bubbletea/v2"
)

type BtListModel struct {
	width   int
	height  int
	focused bool
	// devices []BluetoothDevice // TODO: implement this
	cursor int
}

func (b BtListModel) Init() tea.Cmd {
	return nil
}

func (b BtListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (b BtListModel) View() tea.View {
	return tea.NewView("hello bt world.")
}
