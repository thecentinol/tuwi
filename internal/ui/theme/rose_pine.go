package theme

import "charm.land/lipgloss/v2"

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
	// BG: lipgloss.Color(base),
	BG: nil,
	FG: lipgloss.Color(text),

	Border:               lipgloss.Color(muted),
	BorderFocused:        lipgloss.Color(foam),
	BorderContent:        lipgloss.Color(subtle),
	BorderContentFocused: lipgloss.Color(foam),

	Success: lipgloss.Color(pine),

	Table: TableTheme{
		HeaderBG: lipgloss.Color(overlay),
		HeaderFG: lipgloss.Color(text),

		SelectedBG: lipgloss.Color(iris),
		SelectedFG: lipgloss.Color(overlay),
	},
	Error: ErrorTheme{
		Text:          lipgloss.Color(love),
		Border:        lipgloss.Color(love),
		BorderContent: lipgloss.Color(love),
	},
	Password: PasswordTheme{
		Border:        lipgloss.Color(foam),
		BorderContent: lipgloss.Color(foam),
	},
	Input: InputTheme{
		TextFocused: lipgloss.Color(text),
		TextBlurred: lipgloss.Color(muted),

		PlaceholderFocused: lipgloss.Color(subtle),
		PlaceholderBlurred: lipgloss.Color(muted),

		PromptFocused: lipgloss.Color(subtle),
		PromptBlurred: lipgloss.Color(muted),

		Cursor: lipgloss.Color(foam),
	},
	Help: HelpTheme{
		Ellipsis:       nil,
		ShortKey:       lipgloss.Color(gold),
		ShortDesc:      lipgloss.Color(muted),
		ShortSeparator: lipgloss.Color(muted),
		FullKey:        nil,
		FullDesc:       nil,
		FullSeparator:  nil,
	},
}

var RosePineMoon = Theme{
	// BG: lipgloss.Color(base),
	BG: nil,
	FG: lipgloss.Color(moonText),

	Border:               lipgloss.Color(moonMuted),
	BorderFocused:        lipgloss.Color(moonFoam),
	BorderContent:        lipgloss.Color(moonSubtle),
	BorderContentFocused: lipgloss.Color(moonFoam),

	Success: lipgloss.Color(moonPine),

	Table: TableTheme{
		HeaderBG: lipgloss.Color(moonOverlay),
		HeaderFG: lipgloss.Color(moonText),

		SelectedBG: lipgloss.Color(moonIris),
		SelectedFG: lipgloss.Color(moonOverlay),
	},
	Error: ErrorTheme{
		Text:          lipgloss.Color(moonLove),
		Border:        lipgloss.Color(moonLove),
		BorderContent: lipgloss.Color(moonLove),
	},
	Password: PasswordTheme{
		Border:        lipgloss.Color(moonFoam),
		BorderContent: lipgloss.Color(moonFoam),
	},
	Input: InputTheme{
		TextFocused: lipgloss.Color(moonText),
		TextBlurred: lipgloss.Color(moonMuted),

		PlaceholderFocused: lipgloss.Color(moonSubtle),
		PlaceholderBlurred: lipgloss.Color(moonMuted),

		PromptFocused: lipgloss.Color(moonSubtle),
		PromptBlurred: lipgloss.Color(moonMuted),

		Cursor: lipgloss.Color(moonFoam),
	},
	Help: HelpTheme{
		Ellipsis:       nil,
		ShortKey:       lipgloss.Color(moonGold),
		ShortDesc:      lipgloss.Color(moonSubtle),
		ShortSeparator: lipgloss.Color(moonSubtle),
		FullKey:        nil,
		FullDesc:       nil,
		FullSeparator:  nil,
	},
}
