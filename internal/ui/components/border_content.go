package components

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

type BorderContent struct {
	BorderStyle,
	ContentStyle lipgloss.Style
}

func NewBorderContent(
	border, content color.Color,
	height int,
) BorderContent {
	return BorderContent{
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Height(height),
		ContentStyle: lipgloss.NewStyle().
			Foreground(content).
			PaddingTop(0).
			PaddingBottom(0).
			PaddingLeft(1).
			PaddingRight(1),
	}
}

func (b BorderContent) Render(label, content string, width int) string {
	var (
		border          = b.BorderStyle.GetBorderStyle()
		topBorderStyler = lipgloss.NewStyle().Foreground(b.BorderStyle.GetBorderTopForeground()).Render
		topLeft         = topBorderStyler(border.TopLeft)
		topRight        = topBorderStyler(border.TopRight)

		renderedLabel = b.ContentStyle.Render(label)
	)

	leftWidth := lipgloss.Width(topLeft)
	labelWidth := lipgloss.Width(renderedLabel)
	rightWidth := lipgloss.Width(topRight)
	gap := max(0, width-(leftWidth+labelWidth+rightWidth))

	topBorder := strings.Repeat(border.Top, gap)

	top := topLeft + renderedLabel + topBorderStyler(topBorder) + topRight

	bottomStyle := b.BorderStyle
	bottom := bottomStyle.
		BorderTop(false).
		Width(width).
		Render(content)

	return top + "\n" + bottom
}
