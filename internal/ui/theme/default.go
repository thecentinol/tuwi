package theme

import "charm.land/lipgloss/v2"

var Default = Theme{
	BG:                   nil,
	FG:                   lipgloss.White,
	Border:               lipgloss.Color("238"),
	BorderFocused:        lipgloss.Green,
	BorderContent:        lipgloss.Color("238"),
	BorderContentFocused: lipgloss.Green,
	Success:              lipgloss.BrightGreen,
	Table: TableTheme{
		HeaderBG:   lipgloss.Color("238"),
		HeaderFG:   lipgloss.White,
		SelectedBG: lipgloss.Green,
		SelectedFG: lipgloss.Black,
	},
	Error: ErrorTheme{
		Text:          lipgloss.Red,
		Border:        lipgloss.Green,
		BorderContent: lipgloss.Red,
	},
	Password: PasswordTheme{
		Border:        lipgloss.Green,
		BorderContent: lipgloss.Green,
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
	KeybindsModal: KeybindsModalTheme{
		Border:          lipgloss.Green,
		BorderContent:   lipgloss.Green,
		GroupText:       lipgloss.Green,
		KeybindsText:    lipgloss.Color("#038B5A"),
		DescriptionText: lipgloss.White,
	},
	Help: HelpTheme{
		Ellipsis:       nil,
		ShortKey:       lipgloss.White,
		ShortDesc:      lipgloss.Color("#6e6a86"),
		ShortSeparator: lipgloss.Color("#6e6a86"),
		FullKey:        nil,
		FullDesc:       nil,
		FullSeparator:  nil,
	},
}
