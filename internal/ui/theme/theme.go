package theme

import "image/color"

type Theme struct {
	BG,
	FG,
	Border,
	BorderFocused,
	BorderContent,
	BorderContentFocused,
	Success color.Color

	Table    TableTheme
	Error    ErrorTheme
	Password PasswordTheme
	Input    InputTheme
	Help     HelpTheme
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
