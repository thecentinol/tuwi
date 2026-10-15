package theme

import (
	"errors"
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
	BG HexColor
	FG HexColor

	Border               HexColor
	BorderFocused        HexColor
	BorderContent        HexColor
	BorderContentFocused HexColor

	Success HexColor

	Table         TableTheme
	Error         ErrorTheme
	Password      PasswordTheme
	Input         InputTheme
	KeybindsModal KeybindsModalTheme
	Help          HelpTheme
}

type TableTheme struct {
	HeaderBG   HexColor `toml:"header-bg"`
	HeaderFG   HexColor `toml:"header-fg"`
	CellBG     HexColor `toml:"cell-bg"`
	CellFG     HexColor `toml:"cell-fg"`
	SelectedBG HexColor `toml:"selected-bg"`
	SelectedFG HexColor `toml:"selected-fg"`
}

// error modal
type ErrorTheme struct {
	Text          HexColor `toml:"text"`
	Border        HexColor `toml:"border"`
	BorderContent HexColor `toml:"border-content"`
}

// password modal
type PasswordTheme struct {
	Border        HexColor `toml:"border"`
	BorderContent HexColor `toml:"border-content"`
}

type InputTheme struct {
	TextFocused HexColor `toml:"text-focused"`
	TextBlurred HexColor `toml:"text-muted"`

	PlaceholderFocused HexColor `toml:"placeholder-focused"`
	PlaceholderBlurred HexColor `toml:"placeholder-muted"`

	PromptFocused HexColor `toml:"prompt-focused"`
	PromptBlurred HexColor `toml:"prompt-muted"`

	Cursor HexColor `toml:"cursor"`
}

type KeybindsModalTheme struct {
	Border          HexColor `toml:"border"`
	BorderContent   HexColor `toml:"border-content"`
	GroupText       HexColor `toml:"group-text"`
	KeybindsText    HexColor `toml:"keybinds-text"`
	DescriptionText HexColor `toml:"description-text"`
}

// the help view
type HelpTheme struct {
	Ellipsis       HexColor `toml:"ellipsis"`
	ShortKey       HexColor `toml:"short-key"`
	ShortDesc      HexColor `toml:"short-description"`
	ShortSeparator HexColor `toml:"short-separator"`
	FullKey        HexColor `toml:"full-key"`
	FullDesc       HexColor `toml:"full-description"`
	FullSeparator  HexColor `toml:"full-separator"`
}
