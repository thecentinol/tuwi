package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

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

// If an incorrect theme is provided but the rest of the config is valid, the default theme must be
// set. Irrespective of theme validity, if the provided keybinds are invalid then an error must be
// returned with a nil config - same vice versa.
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
			name: "invalid theme falls back to default",
			content: `
			theme = "invalid"
			`,
			wantTheme: theme.Default,
		},
		{
			name: "valid theme",
			content: `
			theme = "rose-pine"
			`,
			wantTheme: theme.RosePine,
		},
		{
			name: "valid theme with invalid keybind field",
			content: `
			theme = "rose-pine"
			[keybinds]
			auto_connecttt = ["a"]
			`,
			wantErr:   true,
			wantTheme: theme.Default,
		},
		{
			name: "valid theme with malformed keybind",
			content: `
			theme = "rose-pine"
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
