package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type HexColor struct {
	color.Color
}

func (h *HexColor) UnmarshalText(text []byte) error {
	parsed, err := parseHexColor(string(text))
	if err != nil {
		return err
	}

	h.Color = parsed
	return nil
}

func parseHexColor(hex string) (color.Color, error) {
	if hex == "" {
		return nil, errors.New("parseHexColor: nil hex string")
	}
	return lipgloss.Color(hex), nil
}

func (h HexColor) RGBA() (r, g, b, a uint32) {
	if h.Color == nil {
		return 0, 0, 0, 0
	}

	return h.Color.RGBA()
}

type Theme struct {
	BG,
	FG,

	Border,
	BorderFocused,
	BorderContent,
	BorderContentFocused,

	Success color.Color

	Table         TableTheme
	Error         ErrorTheme
	Password      PasswordTheme
	Input         InputTheme
	KeybindsModal KeybindsModalTheme
	Help          HelpTheme
}

type TableTheme struct {
	HeaderBG,
	HeaderFG,
	CellBG,
	CellFG,
	SelectedBG,
	SelectedFG color.Color
}

// error modal
type ErrorTheme struct {
	Text,
	Border,
	BorderContent color.Color
}

// password modal
type PasswordTheme struct {
	Border,
	BorderContent color.Color
}

type InputTheme struct {
	TextFocused,
	TextBlurred,

	PlaceholderFocused,
	PlaceholderBlurred,

	PromptFocused,
	PromptBlurred,

	Cursor color.Color
}

type KeybindsModalTheme struct {
	Border,
	BorderContent,
	GroupText,
	KeybindsText,
	DescriptionText color.Color
}

// the help view
type HelpTheme struct {
	Ellipsis,
	ShortKey,
	ShortDesc,
	ShortSeparator,
	FullKey,
	FullDesc,
	FullSeparator color.Color
}
