package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	kb "github.com/thecentinol/tuwi/internal/keybindings"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

// temporary config for decoding into, see loadFile()
type config struct {
	// NOTE: unexported fields can't be populated by toml decoder
	Theme    customTheme `toml:"theme"`
	Keybinds keybinds    `toml:"keybinds"`
}

type customTheme struct {
	Preset string         `toml:"preset"`
	BG     theme.HexColor `toml:"bg"`
	FG     theme.HexColor `toml:"fg"`

	Border               theme.HexColor `toml:"border"`
	BorderFocused        theme.HexColor `toml:"border-focused"`
	BorderContent        theme.HexColor `toml:"border-content"`
	BorderContentFocused theme.HexColor `toml:"border-content-focused"`

	Success theme.HexColor `toml:"success"`

	Table         theme.TableTheme         `toml:"table"`
	Error         theme.ErrorTheme         `toml:"error-modal"`
	Password      theme.PasswordTheme      `toml:"password-modal"`
	Input         theme.InputTheme         `toml:"input"`
	KeybindsModal theme.KeybindsModalTheme `toml:"keybinds-modal"`
	Help          theme.HelpTheme          `toml:"help-view"`
}

type keybinds struct {
	ConnectSaved             []string `toml:"connect_saved"`
	Disconnect               []string `toml:"disconnect"`
	Forget                   []string `toml:"forget"`
	AutoConnect              []string `toml:"auto_connect"`
	Scan                     []string `toml:"scan"`
	ConnectAvailable         []string `toml:"connect_available"`
	Up                       []string `toml:"up"`
	Down                     []string `toml:"down"`
	GoToTop                  []string `toml:"go_to_top"`
	GoToBottom               []string `toml:"go_to_bottom"`
	DismissError             []string `toml:"dismiss_error"`
	OpenKeybindsModal        []string `toml:"open_keybinds_modal"`
	TogglePasswordVisibility []string `toml:"toggle_password_visibility"`
	SubmitPassword           []string `toml:"submit_password"`
	Cancel                   []string `toml:"cancel"`
	Quit                     []string `toml:"quit"`
}

// the config thats used throughout the app
type AppConfig struct {
	Theme    theme.Theme
	Keybinds kb.Keybindings
}

// applies non-nil colors to the *theme.Theme for the following fields exclusively:
// BG, FG, Border, BorderFocused, BorderContent, BorderContentFocused, Success
func (c *customTheme) applyTo(t *theme.Theme) {
	if c.BG.Color != nil {
		t.BG = c.BG
	}
	if c.FG.Color != nil {
		t.FG = c.FG
	}
	if c.Border.Color != nil {
		t.Border = c.Border
	}
	if c.BorderFocused.Color != nil {
		t.BorderFocused = c.BorderFocused
	}
	if c.BorderContent.Color != nil {
		t.BorderContent = c.BorderContent
	}
	if c.BorderContentFocused.Color != nil {
		t.BorderContentFocused = c.BorderContentFocused
	}
	if c.Success.Color != nil {
		t.Success = c.Success
	}
}

func Load() (*AppConfig, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "tuwi", "config.toml")
	cfg, err := loadFromPath(path)
	return cfg, err
}

func loadFromPath(path string) (*AppConfig, error) {
	cfg, err := loadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	return applyCfg(cfg), nil
}

func loadFile(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg config

	decoder := toml.NewDecoder(file)
	decoder.DisallowUnknownFields()
	decoder.EnableUnmarshalerInterface()

	var derr *toml.DecodeError

	errDec := decoder.Decode(&cfg)
	if errors.As(errDec, &derr) {
		row, col := derr.Position()
		return nil, fmt.Errorf("error %s occurred at row %v column %v", derr.String(), row, col)
	}
	validThemes := map[string]struct{}{
		"default":        {},
		"rose-pine":      {},
		"rose-pine-moon": {},
		"custom":         {},
	}
	if _, valid := validThemes[cfg.Theme.Preset]; !valid && cfg.Theme.Preset != "" {
		return nil, fmt.Errorf("invalid theme preset")
	}

	return &cfg, nil
}

func defaultConfig() *AppConfig {
	return &AppConfig{
		Theme:    theme.Default,
		Keybinds: kb.Default,
	}
}

func applyCfg(cfg *config) *AppConfig {
	appCfg := defaultConfig()

	t := applyTheme(&cfg.Theme)
	appCfg.Theme = *t

	applyKeys := func(dst *kb.Binding, keys []string) {
		if len(keys) > 0 && keys[0] != "" {
			dst.Keys = keys
			dst.HelpKey = keys[0]
		}
	}

	applyKeys(&appCfg.Keybinds.WifiConnectSaved, cfg.Keybinds.ConnectSaved)
	applyKeys(&appCfg.Keybinds.WifiDisconnect, cfg.Keybinds.Disconnect)
	applyKeys(&appCfg.Keybinds.WifiForget, cfg.Keybinds.Forget)
	applyKeys(&appCfg.Keybinds.WifiAutoConnect, cfg.Keybinds.AutoConnect)
	applyKeys(&appCfg.Keybinds.WifiScan, cfg.Keybinds.Scan)
	applyKeys(&appCfg.Keybinds.WifiConnectAvailable, cfg.Keybinds.ConnectAvailable)
	applyKeys(&appCfg.Keybinds.Up, cfg.Keybinds.Up)
	applyKeys(&appCfg.Keybinds.Down, cfg.Keybinds.Down)
	applyKeys(&appCfg.Keybinds.GotoTop, cfg.Keybinds.GoToTop)
	applyKeys(&appCfg.Keybinds.GotoBottom, cfg.Keybinds.GoToBottom)
	applyKeys(&appCfg.Keybinds.ErrorDismiss, cfg.Keybinds.DismissError)
	applyKeys(&appCfg.Keybinds.OpenKeybindingsModal, cfg.Keybinds.OpenKeybindsModal)
	applyKeys(&appCfg.Keybinds.PasswordVisibility, cfg.Keybinds.TogglePasswordVisibility)
	applyKeys(&appCfg.Keybinds.PasswordSubmit, cfg.Keybinds.SubmitPassword)
	applyKeys(&appCfg.Keybinds.Cancel, cfg.Keybinds.Cancel)
	applyKeys(&appCfg.Keybinds.Quit, cfg.Keybinds.Quit)

	return appCfg
}

func applyTheme(ct *customTheme) *theme.Theme {
	var t theme.Theme
	switch ct.Preset {
	case "rose-pine":
		t = theme.RosePine
	case "rose-pine-moon":
		t = theme.RosePineMoon
	case "custom":
		// allows full custom theming
	default:
		t = theme.Default
	}

	ct.applyTo(&t)
	ct.Table.ApplyTo(&t)
	ct.Error.ApplyTo(&t)
	ct.Password.ApplyTo(&t)
	ct.Input.ApplyTo(&t)
	ct.KeybindsModal.ApplyTo(&t)
	ct.Help.ApplyTo(&t)

	return &t
}
