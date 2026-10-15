package theme

import "charm.land/lipgloss/v2"

var Default = Theme{
	FG:                   HexColor{Color: lipgloss.White},
	Border:               HexColor{Color: lipgloss.Color("238")},
	BorderFocused:        HexColor{Color: lipgloss.Green},
	BorderContent:        HexColor{Color: lipgloss.Color("238")},
	BorderContentFocused: HexColor{Color: lipgloss.Green},
	Success:              HexColor{Color: lipgloss.BrightGreen},
	Table: TableTheme{
		HeaderBG:   HexColor{Color: lipgloss.Color("238")},
		HeaderFG:   HexColor{Color: lipgloss.White},
		SelectedBG: HexColor{Color: lipgloss.Green},
		SelectedFG: HexColor{Color: lipgloss.Black},
	},
	Error: ErrorTheme{
		Text:          HexColor{Color: lipgloss.Red},
		Border:        HexColor{Color: lipgloss.Green},
		BorderContent: HexColor{Color: lipgloss.Red},
	},
	Password: PasswordTheme{
		Border:        HexColor{Color: lipgloss.Green},
		BorderContent: HexColor{Color: lipgloss.Green},
	},
	Input: InputTheme{
		TextFocused:        HexColor{Color: lipgloss.White},
		TextBlurred:        HexColor{Color: lipgloss.Color("238")},
		PlaceholderFocused: HexColor{Color: lipgloss.Color("238")},
		PlaceholderBlurred: HexColor{Color: lipgloss.Color("238")},
		PromptFocused:      HexColor{Color: lipgloss.Color("238")},
		PromptBlurred:      HexColor{Color: lipgloss.Color("238")},
		Cursor:             HexColor{Color: lipgloss.White},
	},
	KeybindsModal: KeybindsModalTheme{
		Border:          HexColor{Color: lipgloss.Green},
		BorderContent:   HexColor{Color: lipgloss.Green},
		GroupText:       HexColor{Color: lipgloss.Green},
		KeybindsText:    HexColor{Color: lipgloss.Color("#038B5A")},
		DescriptionText: HexColor{Color: lipgloss.White},
	},
	Help: HelpTheme{
		ShortKey:       HexColor{Color: lipgloss.White},
		ShortDesc:      HexColor{Color: lipgloss.Color("#6e6a86")},
		ShortSeparator: HexColor{Color: lipgloss.Color("#6e6a86")},
	},
}
