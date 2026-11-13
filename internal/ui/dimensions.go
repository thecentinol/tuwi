package ui

import (
	"charm.land/lipgloss/v2"

	comp "github.com/thecentinol/tuwi/internal/ui/components"
	apphelp "github.com/thecentinol/tuwi/internal/ui/help"
	wifiui "github.com/thecentinol/tuwi/internal/ui/wifi"
)

func (m *Model) setDimensions() {
	helpViewHeight := 1

	setWifiDimensions(&m.wifiSaved, &m.wifiAvailable, (m.width*2)/3, m.height-helpViewHeight)
	setHelpMenuDimensions(&m.HelpMenu, m.width, m.height)
	setPasswordModalDimensionsAndPos(&m.passwordModal, m.width, m.height, m.width/3)
	setErrorModalDimensionsAndPos(&m.errorModal, m.width, m.height)
	m.infoPane.SetDimensions(m.width/3, m.height-helpViewHeight)
}

func setWifiDimensions(
	saved *wifiui.SavedModel,
	available *wifiui.AvailableModel,
	width, height int,
) {
	savedHeight := height / 2
	availableHeight := height - savedHeight

	saved.SetDimensions(width, savedHeight)
	available.SetDimensions(width, availableHeight)
}

func setHelpMenuDimensions(m *apphelp.Model, parentWidth, parentHeight int) {
	// 2 for borders and another 2 for a little extra
	// because it looks better
	width := parentWidth - 4
	height := parentHeight - 4

	// TODO: use consts instead of magic numbers
	if parentHeight > 35 {
		height = 28
	}
	if parentWidth > 95 {
		width = 90
	}

	m.SetDimensions(width, height)

	content := m.View().Content
	overlayX := (parentWidth / 2) - (lipgloss.Width(content) / 2)
	overlayY := (parentHeight / 2) - (lipgloss.Height(content) / 2)
	m.SetOverlayPosition(overlayX, overlayY)
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
