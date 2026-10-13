package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/thecentinol/tuwi/internal/ui/theme"
)

type Config struct {
	Theme string `toml:"theme"`
}

type AppConfig struct {
	Theme theme.Theme
}

func Load() (*AppConfig, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "tuwi", "config.toml")

	var cfg Config
	appCfg := &AppConfig{
		Theme: theme.Default,
	}

	file, err := os.Open(path)
	if errors.Is(err, os.ErrPermission) {
		return appCfg, err
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	switch cfg.Theme {
	case "rose-pine":
		appCfg.Theme = theme.RosePine
	}

	return appCfg, nil
}
