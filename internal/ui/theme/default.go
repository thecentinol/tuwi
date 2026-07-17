package theme

import "charm.land/lipgloss/v2"

var Default = Theme{
	BG:            nil,
	FG:            lipgloss.White,
	Border:        lipgloss.Color("238"),
	BorderFocused: lipgloss.Green,
	Success:       lipgloss.BrightGreen,
	Table: TableTheme{
		HeaderBG:   lipgloss.Color("238"),
		HeaderFG:   lipgloss.White,
		SelectedBG: lipgloss.Green,
		SelectedFG: lipgloss.Black,
	},
	Error: ErrorTheme{
		Text:   lipgloss.Red,
		Border: lipgloss.Green,
	},
	Input: InputTheme{
		TextFocused:        lipgloss.White,
		TextBlurred:        lipgloss.Color("238"),
		PlaceholderFocused: lipgloss.Color("238"),
		PlaceholderBlurred: lipgloss.Color("238"),
		PromptFocused:      lipgloss.Color("238"),
		PromptBlurred:      lipgloss.Color("238"),
		Cursor:             lipgloss.White,
	},
	Help: HelpTheme{
		Ellipsis:       nil,
		ShortKey:       lipgloss.White,
		ShortDesc:      lipgloss.Color("238"),
		ShortSeparator: lipgloss.Color("238"),
		FullKey:        nil,
		FullDesc:       nil,
		FullSeparator:  nil,
	},
}
