package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	kb "github.com/thecentinol/tuwi/internal/keybindings"
	"github.com/thecentinol/tuwi/internal/ui/theme"
)

// temporary config for decoding into, see Load()
type config struct {
	// NOTE: unexported fields can't be populated by toml decoder
	Theme    string   `toml:"theme"`
	Keybinds keybinds `toml:"keybinds"`
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
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
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

	switch cfg.Theme {
	case "rose-pine":
		appCfg.Theme = theme.RosePine
	case "rose-pine-moon":
		appCfg.Theme = theme.RosePineMoon
	}

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
