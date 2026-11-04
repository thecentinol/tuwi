package config

import (
	"fmt"
	"image/color"
	"regexp"

	"charm.land/lipgloss/v2"
)

var Themes = map[string]Theme{
	"base":           DefaultTheme,
	"rose-pine":      RosePine,
	"rose-pine-moon": RosePineMoon,
}

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
	matched, err := regexp.MatchString(`^#[0-9a-fA-F]{6}$`, hex)
	if err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("invalid hex code")
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

	Table    TableTheme
	Error    ErrorTheme
	Password PasswordTheme
	Input    InputTheme
	HelpMenu HelpMenuTheme
	HelpView HelpViewTheme
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

type HelpMenuTheme struct {
	Border          HexColor `toml:"border"`
	BorderContent   HexColor `toml:"border-content"`
	GroupText       HexColor `toml:"group-text"`
	KeybindsText    HexColor `toml:"keybinds-text"`
	DescriptionText HexColor `toml:"description-text"`
}

// the help view
type HelpViewTheme struct {
	Ellipsis       HexColor `toml:"ellipsis"`
	ShortKey       HexColor `toml:"short-key"`
	ShortDesc      HexColor `toml:"short-description"`
	ShortSeparator HexColor `toml:"short-separator"`
	FullKey        HexColor `toml:"full-key"`
	FullDesc       HexColor `toml:"full-description"`
	FullSeparator  HexColor `toml:"full-separator"`
}

// apply non-nil colors to the *Theme for the TableTheme field
func (tt *TableTheme) ApplyTo(t *Theme) {
	if tt.HeaderBG.Color != nil {
		t.Table.HeaderBG = tt.HeaderBG
	}
	if tt.HeaderFG.Color != nil {
		t.Table.HeaderFG = tt.HeaderFG
	}
	if tt.CellBG.Color != nil {
		t.Table.CellBG = tt.CellBG
	}
	if tt.CellFG.Color != nil {
		t.Table.CellFG = tt.CellFG
	}
	if tt.SelectedBG.Color != nil {
		t.Table.SelectedBG = tt.SelectedBG
	}
	if tt.SelectedFG.Color != nil {
		t.Table.SelectedFG = tt.SelectedFG
	}
}

// same as ApplyTo on TableTheme, but for the ErrorTheme
func (e *ErrorTheme) ApplyTo(t *Theme) {
	if e.Text.Color != nil {
		t.Error.Text = e.Text
	}
	if e.Border.Color != nil {
		t.Error.Border = e.Border
	}
	if e.BorderContent.Color != nil {
		t.Error.BorderContent = e.BorderContent
	}
}

// same as ApplyTo on TableTheme, but for the PasswordTheme
func (p *PasswordTheme) ApplyTo(t *Theme) {
	if p.Border.Color != nil {
		t.Password.Border = p.Border
	}
	if p.BorderContent.Color != nil {
		t.Password.BorderContent = p.BorderContent
	}
}

// same as ApplyTo on TableTheme, but for the InputTheme
func (i *InputTheme) ApplyTo(t *Theme) {
	if i.TextFocused.Color != nil {
		t.Input.TextFocused = i.TextFocused
	}
	if i.TextBlurred.Color != nil {
		t.Input.TextBlurred = i.TextBlurred
	}
	if i.PlaceholderFocused.Color != nil {
		t.Input.PlaceholderFocused = i.PlaceholderFocused
	}
	if i.PlaceholderBlurred.Color != nil {
		t.Input.PlaceholderBlurred = i.PlaceholderBlurred
	}
	if i.PromptFocused.Color != nil {
		t.Input.PromptFocused = i.PromptFocused
	}
	if i.PromptBlurred.Color != nil {
		t.Input.PromptBlurred = i.PromptBlurred
	}
	if i.Cursor.Color != nil {
		t.Input.Cursor = i.Cursor
	}
}

// same as ApplyTo on TableTheme, but for the HelpMenuTheme
func (k *HelpMenuTheme) ApplyTo(t *Theme) {
	if k.Border.Color != nil {
		t.HelpMenu.Border = k.Border
	}
	if k.BorderContent.Color != nil {
		t.HelpMenu.BorderContent = k.BorderContent
	}
	if k.GroupText.Color != nil {
		t.HelpMenu.GroupText = k.GroupText
	}
	if k.KeybindsText.Color != nil {
		t.HelpMenu.KeybindsText = k.KeybindsText
	}
	if k.DescriptionText.Color != nil {
		t.HelpMenu.DescriptionText = k.DescriptionText
	}
}

// same as ApplyTo on TableTheme, but for the HelpTheme
func (h *HelpViewTheme) ApplyTo(t *Theme) {
	if h.Ellipsis.Color != nil {
		t.HelpView.Ellipsis = h.Ellipsis
	}
	if h.ShortKey.Color != nil {
		t.HelpView.ShortKey = h.ShortKey
	}
	if h.ShortDesc.Color != nil {
		t.HelpView.ShortDesc = h.ShortDesc
	}
	if h.ShortSeparator.Color != nil {
		t.HelpView.ShortSeparator = h.ShortSeparator
	}
	if h.FullKey.Color != nil {
		t.HelpView.FullKey = h.FullKey
	}
	if h.FullDesc.Color != nil {
		t.HelpView.FullDesc = h.FullDesc
	}
	if h.FullSeparator.Color != nil {
		t.HelpView.FullSeparator = h.FullSeparator
	}
}
