package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

// Init creates the log dir and file and
// sets the output of all logs to that file
func Init() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logDir := filepath.Join(home, ".local", "state", "tuwi")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, "tuwi.log")

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(
		slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	)

	return file, nil
}
