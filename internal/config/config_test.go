package config

import (
	"errors"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/pelletier/go-toml/v2"

	"github.com/thecentinol/tuwi/internal/ui/theme"
)

// test if decoded config keybinds are loaded correctly
func TestDecodeKeybinds(t *testing.T) {
	// the default keybinds that users can configure
	cfgInput := `
    [keybinds]
    connect_saved = ["enter"]
    disconnect = ["d"]
    forget = ["f"]
    auto_connect = ["a"]
    scan = ["s"]
    connect_available = ["enter"]
    up = ["k", "up"]
    down = ["j", "down"]
    go_to_top = ["t"]
    go_to_bottom = ["b"]
    dismiss_error = ["esc"]
    cancel = ["esc"]
    quit = ["q", "ctrl+c"]
    `

	var cfg config
	err := toml.NewDecoder(strings.NewReader(cfgInput)).
		DisallowUnknownFields().
		Decode(&cfg)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(cfg.Keybinds.ConnectSaved, []string{"enter"}) {
		t.Errorf("unexpected connect_saved keybinds: %v", cfg.Keybinds.ConnectSaved)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Disconnect, []string{"d"}) {
		t.Errorf("unexpected disconnect keybinds: %v", cfg.Keybinds.Disconnect)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Forget, []string{"f"}) {
		t.Errorf("unexpected forget keybinds: %v", cfg.Keybinds.Forget)
	}
	if !reflect.DeepEqual(cfg.Keybinds.AutoConnect, []string{"a"}) {
		t.Errorf("unexpected auto_connect keybinds: %v", cfg.Keybinds.AutoConnect)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Scan, []string{"s"}) {
		t.Errorf("unexpected scan keybinds: %v", cfg.Keybinds.Scan)
	}
	if !reflect.DeepEqual(cfg.Keybinds.ConnectAvailable, []string{"enter"}) {
		t.Errorf("unexpected connect_available keybinds: %v", cfg.Keybinds.ConnectAvailable)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Up, []string{"k", "up"}) {
		t.Errorf("unexpected up keybinds: %v", cfg.Keybinds.Up)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Down, []string{"j", "down"}) {
		t.Errorf("unexpected down keybinds: %v", cfg.Keybinds.Down)
	}
	if !reflect.DeepEqual(cfg.Keybinds.GoToTop, []string{"t"}) {
		t.Errorf("unexpected go_to_top keybinds: %v", cfg.Keybinds.GoToTop)
	}
	if !reflect.DeepEqual(cfg.Keybinds.GoToBottom, []string{"b"}) {
		t.Errorf("unexpected go_to_bottom keybinds: %v", cfg.Keybinds.GoToBottom)
	}
	if !reflect.DeepEqual(cfg.Keybinds.DismissError, []string{"esc"}) {
		t.Errorf("unexpected dismiss_error keybinds: %v", cfg.Keybinds.DismissError)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Cancel, []string{"esc"}) {
		t.Errorf("unexpected cancel keybinds: %v", cfg.Keybinds.Cancel)
	}
	if !reflect.DeepEqual(cfg.Keybinds.Quit, []string{"q", "ctrl+c"}) {
		t.Errorf("unexpected quit keybinds: %v", cfg.Keybinds.Quit)
	}
}

func TestLoadMissingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	_, err := loadFile(path)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

// for a missing config file/dir loadFromPath() must always return defaultConfig()
// and never return an err
func TestMissingConfigReturnsDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")

	got, err := loadFromPath(path)
	want := defaultConfig()

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected default config, got %v with error = %v", got, err)
	}
}

// If an incorrect theme is provided but the rest of the config is valid, the default theme must
// NOT be set and an error must be returned. With a valid theme if the provided keybinds are
// invalid then an error must be returned with a nil config - same vice versa.
func TestLoadFromPathReturns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   string
		wantErr   bool
		wantTheme theme.Theme
	}{
		{
			name: "missing double quote",
			content: `
			[keybinds]
			connect_saved = ["enter]
			`,
			wantErr:   true,
			wantTheme: theme.Default,
		},
		{
			name: "incorrect keybind field name",
			content: `
			[keybinds]
			connectttt_saved = ["enter"]
			disconnect = ["d"]
			forget = ["f"]
			auto_connect = ["a"]
			scan = ["s"]
			connect_available = ["enter"]
			up = ["k", "up"]
			down = ["j", "down"]
			go_to_top = ["t"]
			go_to_bottom = ["b"]
			dismiss_error = ["esc"]
			cancel = ["esc"]
			quit = ["q", "ctrl+c"]
			`,
			wantErr:   true,
			wantTheme: theme.Default,
		},
		{
			name: "empty keybind string",
			content: `
			[keybinds]
			disconnect = [""]
			`,
			wantTheme: theme.Default,
		},
		{
			name: "valid theme",
			content: `
            [theme]
            preset = "rose-pine"
			`,
			wantTheme: theme.RosePine,
		},
		{
			name: "valid theme with invalid keybind field",
			content: `
            [theme]
            preset = "rose-pine"
			[keybinds]
			auto_connecttt = ["a"]
			`,
			wantErr:   true,
			wantTheme: theme.Default,
		},
		{
			name: "valid theme with malformed keybind",
			content: `
            [theme]
            preset = "rose-pine"
			[keybinds]
			go_to_top = ["t]
			`,
			wantErr:   true,
			wantTheme: theme.Default,
		},
		{
			name: "empty dismiss error keybind",
			content: `
			[keybinds]
			dismiss_error = [""]
			`,
			wantTheme: theme.Default,
		},
		{
			name: "valid theme without keybinds",
			content: `
            [theme]
            preset = "rose-pine-moon"
			`,
			wantTheme: theme.RosePineMoon,
		},
		{
			name: "valid theme valid keybinds",
			content: `
            [theme]
            preset = "rose-pine-moon"
			[keybinds]
			connect_saved = ["enter"]
			disconnect = ["d"]
			forget = ["f"]
			auto_connect = ["a"]
			scan = ["s"]
			connect_available = ["enter"]
			up = ["k", "up"]
			down = ["j", "down"]
			go_to_top = ["t"]
			go_to_bottom = ["b"]
			dismiss_error = ["esc"]
			cancel = ["esc"]
			quit = ["q", "ctrl+c"]
			`,
			wantTheme: theme.RosePineMoon,
		},
		{
			name: "default theme",
			content: `
            [theme]
            preset = "default"
			`,
			wantErr:   false,
			wantTheme: theme.Default,
		},
		{
			name: "defaultt theme",
			content: `
            [theme]
            preset = "defaultt"
			`,
			wantErr: true,
		},
		{
			name: "invalid theme",
			content: `
            [theme]
            preset = "invalid"
			`,
			wantErr: true,
		},
		{
			name: "rosee-pine theme",
			content: `
            [theme]
            preset = "rosee-pine"
			`,
			wantErr: true,
		},
		{
			name: "rose-pine-moon theme",
			content: `
            [theme]
            preset = "rose-pine-moon"
			`,
			wantErr:   false,
			wantTheme: theme.RosePineMoon,
		},
		{
			name: "rose-pine-sun",
			content: `
            [theme]
            preset = "rose-pine-sun"
			`,
			wantErr: true,
		},
		{
			name: "random-theme-name",
			content: `
            [theme]
            preset = "random-theme-name"
			`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatal(err)
			}

			cfg, err := loadFromPath(path)

			if (err != nil) != tc.wantErr {
				t.Fatalf("loadFromPath() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr {
				if cfg != nil {
					t.Errorf("loadFromPath() config = %v, want nil on error", cfg)
				}
				return
			}

			if cfg == nil {
				t.Fatal("loadFromPath() unexpectedly returned nil config")
			}

			if cfg.Theme != tc.wantTheme {
				t.Fatalf("loadFromPath() theme = %v, want %v", cfg.Theme, tc.wantTheme)
			}
		})
	}
}

// test that the default theme's colors are applied, including the nested sections' colors
func TestApplyTheme_PresetDefault(t *testing.T) {
	cfg := &customTheme{
		Preset: "default",
	}

	got := applyTheme(cfg)

	if got.BG.Color != theme.Default.BG.Color {
		t.Errorf(
			"BG = %v, want default BG %v",
			got.BG.Color, theme.Default.BG.Color,
		)
	}
	if got.Border.Color != theme.Default.Border.Color {
		t.Errorf(
			"Border = %v, want default Border %v",
			got.Border.Color, theme.Default.Border.Color,
		)
	}
	if got.Error.Text.Color != theme.Default.Error.Text.Color {
		t.Errorf(
			"ErrorModal Text = %v, want default ErrorModal Text %v",
			got.Error.Text.Color, theme.Default.Error.Text.Color,
		)
	}
	if got.Password.BorderContent.Color != theme.Default.Password.BorderContent.Color {
		t.Errorf(
			"PasswordModal BorderContent = %v, want default PasswordModal BorderContent %v",
			got.Password.BorderContent.Color, theme.Default.Password.BorderContent.Color,
		)
	}
}

// test that the RosePine theme's colors are applied, including the nested sections' colors
func TestApplyTheme_PresetRosePine(t *testing.T) {
	cfg := &customTheme{
		Preset: "rose-pine",
	}

	got := applyTheme(cfg)

	if got.BG.Color != theme.RosePine.BG.Color {
		t.Errorf(
			"BG = %v, want RosePine BG %v",
			got.BG.Color, theme.RosePine.BG.Color,
		)
	}
	if got.Border.Color != theme.RosePine.Border.Color {
		t.Errorf(
			"Border = %v, want RosePine Border %v",
			got.Border.Color, theme.RosePine.Border.Color,
		)
	}
	if got.Error.Text.Color != theme.RosePine.Error.Text.Color {
		t.Errorf(
			"ErrorModal Text = %v, want RosePine ErrorModal Text %v",
			got.Error.Text.Color, theme.RosePine.Error.Text.Color,
		)
	}
	if got.Password.BorderContent.Color != theme.RosePine.Password.BorderContent.Color {
		t.Errorf(
			"PasswordModal BorderContent = %v, want RosePine PasswordModal BorderContent %v",
			got.Password.BorderContent.Color, theme.RosePine.Password.BorderContent.Color,
		)
	}
}

// test that the RosePineMoon theme's colors are applied, including the nested sections' colors
func TestApplyTheme_PresetRosePineMoon(t *testing.T) {
	cfg := &customTheme{
		Preset: "rose-pine",
	}

	got := applyTheme(cfg)

	if got.BG.Color != theme.RosePineMoon.BG.Color {
		t.Errorf(
			"BG = %v, want RosePineMoon BG %v",
			got.BG.Color, theme.RosePineMoon.BG.Color,
		)
	}
	if got.Border.Color != theme.RosePineMoon.Border.Color {
		t.Errorf(
			"Border = %v, want RosePineMoon Border %v",
			got.Border.Color, theme.RosePineMoon.Border.Color,
		)
	}
	if got.Error.Text.Color != theme.RosePineMoon.Error.Text.Color {
		t.Errorf(
			"ErrorModal Text = %v, want RosePineMoon ErrorModal Text %v",
			got.Error.Text.Color, theme.RosePineMoon.Error.Text.Color,
		)
	}
	if got.Password.BorderContent.Color != theme.RosePineMoon.Password.BorderContent.Color {
		t.Errorf(
			"PasswordModal BorderContent = %v, want RosePineMoon PasswordModal BorderContent %v",
			got.Password.BorderContent.Color, theme.RosePineMoon.Password.BorderContent.Color,
		)
	}
}

// test that the default theme's colors are correctly overridden when customized
func TestApplyTheme_PartialOverride(t *testing.T) {
	cfg := &customTheme{
		Preset: "default",
		BG:     theme.HexColor{Color: lipgloss.Black},
	}

	// set nested sections' colors using lipgloss.<color_name>, color.RGBA{} and hex values
	cfg.Table.HeaderFG = theme.HexColor{Color: lipgloss.White}
	cfg.Input.Cursor = theme.HexColor{Color: color.RGBA{R: 208, G: 12, B: 247}}
	cfg.Help.ShortSeparator = theme.HexColor{Color: lipgloss.Color("#d00cf7")}

	got := applyTheme(cfg)

	if got.FG.Color != theme.Default.FG.Color {
		t.Errorf("FG = %v, want default preset FG %v", got.FG.Color, theme.Default.FG.Color)
	}
	if got.Password.Border.Color != theme.Default.Password.Border.Color {
		t.Errorf(
			"Password Border = %v, want default preset Password Border %v",
			got.Password.Border.Color, theme.Default.Password.Border.Color,
		)
	}
	if got.BG.Color != cfg.BG.Color {
		t.Errorf("BG = %v, want custom BG %v", got.BG.Color, cfg.BG.Color)
	}
	if got.Table.HeaderFG.Color != cfg.Table.HeaderFG.Color {
		t.Errorf(
			"Table Header FG = %v, want custom Header FG %v",
			got.Table.HeaderFG.Color, cfg.Table.HeaderFG.Color,
		)
	}
	if got.Input.Cursor.Color != cfg.Input.Cursor.Color {
		t.Errorf(
			"Input Cursor = %v, want custom Input Cursor %v",
			got.Input.Cursor.Color, cfg.Input.Cursor.Color,
		)
	}
	if got.Help.ShortSeparator.Color != cfg.Help.ShortSeparator.Color {
		t.Errorf(
			"Help ShortSeparator = %v, want custom Help ShortSeparator %v",
			got.Help.ShortSeparator.Color, cfg.Help.ShortSeparator.Color,
		)
	}
}

// test that a fully custom theme's colors are applied with no defaults overriding any
func TestApplyTheme_Custom(t *testing.T) {
	cfg := &customTheme{
		Preset: "custom",
		BG:     theme.HexColor{Color: lipgloss.Blue},
		Border: theme.HexColor{Color: lipgloss.BrightMagenta},
	}

	got := applyTheme(cfg)

	if got.BG.Color != cfg.BG.Color {
		t.Errorf(
			"BG = %v, want custom BG %v",
			got.BG.Color, cfg.BG.Color,
		)
	}
	if got.Border.Color != cfg.Border.Color {
		t.Errorf(
			"Border = %v, want custom Border %v",
			got.Border.Color, cfg.Border.Color,
		)
	}
	if got.FG.Color == theme.Default.FG.Color {
		t.Errorf(
			"FG = %v, unexpected default FG %v",
			got.FG.Color, theme.Default.FG.Color,
		)
	}
	if got.BorderContent.Color == theme.Default.BorderContent.Color {
		t.Errorf(
			"BorderContent = %v, unexpected BorderContent %v",
			got.BorderContent.Color, theme.Default.BorderContent.Color,
		)
	}
	if got.Table.HeaderFG.Color == theme.Default.Table.HeaderFG.Color {
		t.Errorf(
			"Table HeaderFG = %v, unexpected Table HeaderFG %v",
			got.Table.HeaderFG.Color, theme.Default.Table.HeaderFG.Color,
		)
	}
}

// when no theming config is provided the default theme should be applied
func TestApplyTheme_NilConfig(t *testing.T) {
	cfg := &customTheme{}

	got := applyTheme(cfg)

	if got.FG.Color != theme.Default.FG.Color {
		t.Errorf(
			"FG = %v, expected default FG %v",
			got.FG.Color, theme.Default.FG.Color,
		)
	}
	if got.Success.Color != theme.Default.Success.Color {
		t.Errorf(
			"Success = %v, expected default Success %v",
			got.Success.Color, theme.Default.Success.Color,
		)
	}
	if got.KeybindsModal.GroupText.Color != theme.Default.KeybindsModal.GroupText.Color {
		t.Errorf(
			"KeybindsModal GroupText = %v, expected default KeybindsModal GroupText %v",
			got.KeybindsModal.GroupText.Color, theme.Default.KeybindsModal.GroupText.Color,
		)
	}
	if got.Help.ShortKey.Color != theme.Default.Help.ShortKey.Color {
		t.Errorf(
			"Help ShortKey = %v, expected default Help ShortKey %v",
			got.Help.ShortKey.Color, theme.Default.Help.ShortKey.Color,
		)
	}
}
