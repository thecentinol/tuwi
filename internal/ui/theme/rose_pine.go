package theme

import "charm.land/lipgloss/v2"

//nolint:goconst // Rose Pine Default and Moon intentionally share certain hex codes
const (
	base          = "#191724"
	surface       = "#1f1d2e"
	overlay       = "#26233a"
	muted         = "#6e6a86"
	subtle        = "#908caa"
	text          = "#e0def4"
	love          = "#eb6f92"
	gold          = "#f6c177"
	rose          = "#ebbcba"
	pine          = "#31748f"
	foam          = "#9ccfd8"
	iris          = "#c4a7e7"
	highlightLow  = "#21202e"
	highlightMed  = "#403d52"
	highlightHigh = "#524f67"

	moonBase          = "#232136"
	moonSurface       = "#2a273f"
	moonOverlay       = "#393552"
	moonMuted         = "#6e6a86"
	moonSubtle        = "#908caa"
	moonText          = "#e0def4"
	moonLove          = "#eb6f92"
	moonGold          = "#f6c177"
	moonRose          = "#ea9a97"
	moonPine          = "#3e8fb0"
	moonFoam          = "#9ccfd8"
	moonIris          = "#c4a7e7"
	moonHighlightLow  = "#2a283e"
	moonHighlightMed  = "#44415a"
	moonHighlightHigh = "#56526e"
)

var RosePine = Theme{
	// BG: HexColor{Color: lipgloss.Color(base)},
	FG: HexColor{Color: lipgloss.Color(text)},

	Border:               HexColor{Color: lipgloss.Color(muted)},
	BorderFocused:        HexColor{Color: lipgloss.Color(foam)},
	BorderContent:        HexColor{Color: lipgloss.Color(subtle)},
	BorderContentFocused: HexColor{Color: lipgloss.Color(foam)},

	Success: HexColor{Color: lipgloss.Color(pine)},

	Table: TableTheme{
		HeaderBG: HexColor{Color: lipgloss.Color(overlay)},
		HeaderFG: HexColor{Color: lipgloss.Color(text)},

		SelectedBG: HexColor{Color: lipgloss.Color(iris)},
		SelectedFG: HexColor{Color: lipgloss.Color(overlay)},
	},
	Error: ErrorTheme{
		Text:          HexColor{Color: lipgloss.Color(love)},
		Border:        HexColor{Color: lipgloss.Color(love)},
		BorderContent: HexColor{Color: lipgloss.Color(love)},
	},
	Password: PasswordTheme{
		Border:        HexColor{Color: lipgloss.Color(foam)},
		BorderContent: HexColor{Color: lipgloss.Color(foam)},
	},
	Input: InputTheme{
		TextFocused: HexColor{Color: lipgloss.Color(text)},
		TextBlurred: HexColor{Color: lipgloss.Color(muted)},

		PlaceholderFocused: HexColor{Color: lipgloss.Color(subtle)},
		PlaceholderBlurred: HexColor{Color: lipgloss.Color(muted)},

		PromptFocused: HexColor{Color: lipgloss.Color(subtle)},
		PromptBlurred: HexColor{Color: lipgloss.Color(muted)},

		Cursor: HexColor{Color: lipgloss.Color(foam)},
	},
	KeybindsModal: KeybindsModalTheme{
		Border:          HexColor{Color: lipgloss.Color(foam)},
		BorderContent:   HexColor{Color: lipgloss.Color(foam)},
		GroupText:       HexColor{Color: lipgloss.Color(foam)},
		KeybindsText:    HexColor{Color: lipgloss.Color(gold)},
		DescriptionText: HexColor{Color: lipgloss.White},
	},
	Help: HelpTheme{
		ShortKey:       HexColor{Color: lipgloss.Color(gold)},
		ShortDesc:      HexColor{Color: lipgloss.Color(muted)},
		ShortSeparator: HexColor{Color: lipgloss.Color(muted)},
	},
}

var RosePineMoon = Theme{
	// BG: HexColor{Color: lipgloss.Color(base)},
	FG: HexColor{Color: lipgloss.Color(moonText)},

	Border:               HexColor{Color: lipgloss.Color(moonMuted)},
	BorderFocused:        HexColor{Color: lipgloss.Color(moonFoam)},
	BorderContent:        HexColor{Color: lipgloss.Color(moonSubtle)},
	BorderContentFocused: HexColor{Color: lipgloss.Color(moonFoam)},

	Success: HexColor{Color: lipgloss.Color(moonPine)},

	Table: TableTheme{
		HeaderBG: HexColor{Color: lipgloss.Color(moonOverlay)},
		HeaderFG: HexColor{Color: lipgloss.Color(moonText)},

		SelectedBG: HexColor{Color: lipgloss.Color(moonIris)},
		SelectedFG: HexColor{Color: lipgloss.Color(moonOverlay)},
	},
	Error: ErrorTheme{
		Text:          HexColor{Color: lipgloss.Color(moonLove)},
		Border:        HexColor{Color: lipgloss.Color(moonLove)},
		BorderContent: HexColor{Color: lipgloss.Color(moonLove)},
	},
	Password: PasswordTheme{
		Border:        HexColor{Color: lipgloss.Color(moonFoam)},
		BorderContent: HexColor{Color: lipgloss.Color(moonFoam)},
	},
	Input: InputTheme{
		TextFocused: HexColor{Color: lipgloss.Color(moonText)},
		TextBlurred: HexColor{Color: lipgloss.Color(moonMuted)},

		PlaceholderFocused: HexColor{Color: lipgloss.Color(moonSubtle)},
		PlaceholderBlurred: HexColor{Color: lipgloss.Color(moonMuted)},

		PromptFocused: HexColor{Color: lipgloss.Color(moonSubtle)},
		PromptBlurred: HexColor{Color: lipgloss.Color(moonMuted)},

		Cursor: HexColor{Color: lipgloss.Color(moonFoam)},
	},
	KeybindsModal: KeybindsModalTheme{
		Border:          HexColor{Color: lipgloss.Color(moonFoam)},
		BorderContent:   HexColor{Color: lipgloss.Color(moonFoam)},
		GroupText:       HexColor{Color: lipgloss.Color(moonFoam)},
		KeybindsText:    HexColor{Color: lipgloss.Color(moonGold)},
		DescriptionText: HexColor{Color: lipgloss.White},
	},
	Help: HelpTheme{
		ShortKey:       HexColor{Color: lipgloss.Color(moonGold)},
		ShortDesc:      HexColor{Color: lipgloss.Color(moonSubtle)},
		ShortSeparator: HexColor{Color: lipgloss.Color(moonSubtle)},
	},
}
