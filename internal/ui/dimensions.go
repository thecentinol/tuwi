package ui

import (
	"charm.land/lipgloss/v2"

	comp "github.com/thecentinol/tuwi/internal/ui/components"
	kbmodal "github.com/thecentinol/tuwi/internal/ui/keybindings_modal"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
)

func (m *Model) setDimensions() {
	helpViewHeight := 2

	setWifiDimensions(&m.wifiSaved, &m.wifiAvailable, m.width/2, m.height-helpViewHeight,
		m.height-helpViewHeight-2)
	setKeybindsModalDimensions(&m.keybindingsModal, m.width, m.height)
	setPasswordModalDimensionsAndPos(&m.passwordModal, m.width, m.height, m.width/3)
	setErrorModalDimensionsAndPos(&m.errorModal, m.width, m.height)
}

func setWifiDimensions(
	saved *wifiui.SavedModel,
	available *wifiui.AvailableModel,
	width, height, tableHeight int,
) {
	saved.SetWidth(width)
	saved.Setheight(height)
	saved.Table.SetWidth(width)
	saved.Table.SetHeight(tableHeight)

	available.SetWidth(width)
	available.Setheight(height)
	available.Table.SetWidth(width)
	available.Table.SetHeight(tableHeight)
}

func setKeybindsModalDimensions(m *kbmodal.Model, parentWidth, parentHeight int) {
	m.SetSize(parentWidth/2, (parentHeight*3)/4)
	m.X = (parentWidth / 2) - (m.Width / 2)
	m.Y = (parentHeight / 2) - (m.Height / 2)
}

// set the dimensions and X/Y position of the password modal
func setPasswordModalDimensionsAndPos(m *comp.PasswordModel, parentWidth, parentHeight,
	modalWidth int,
) {
	iconWidth := 2 // the show/hide indicator icon on the right of the input
	m.Width = modalWidth
	m.Input.SetWidth(modalWidth - iconWidth - 4)

	// get height and width of password modal
	renderedPasswordModal := m.View().Content
	getPasswordModalHeight := lipgloss.Height(renderedPasswordModal)
	getPasswordModalWidth := lipgloss.Width(renderedPasswordModal)

	// calculate X and Y coordinates of password modal
	m.X = (parentWidth / 2) - (getPasswordModalWidth / 2)
	m.Y = (parentHeight / 2) - (getPasswordModalHeight / 2)
}

func setErrorModalDimensionsAndPos(m *comp.ErrorModel, parentWidth, parentHeight int) {
	m.Width = parentWidth / 3
	m.MaxHeight = parentHeight / 5
	m.X = (parentWidth / 2) - (m.Width / 2)
	m.Y = (parentHeight / 2) - (m.Height / 2)
}
