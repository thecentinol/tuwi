package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/logger"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err) //nolint:revive
		os.Exit(1)
	}
}

func run() error {
	// initialize logger
	logger, err := logger.Init()
	if err != nil {
		slog.SetDefault(slog.New(slog.DiscardHandler))
		fmt.Fprintf(os.Stderr, "warning: logging disabled: %v\n", err) //nolint:revive
	} else {
		defer logger.Close()
	}

	client, err := dbus.NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// load config
	cfg, errCfg := config.Load()

	model := ui.NewModel(client, *cfg)

	if !errors.Is(errCfg, os.ErrNotExist) && errCfg != nil {
		model = model.SetStartupError(
			fmt.Errorf("could not load config, using defaults: %w", errCfg),
		)
	}

	p := tea.NewProgram(model)

	cb := nm.SignalCallbacks{
		OnAccessPointAdded: func(ap nm.AccessPointProperties) {
			p.Send(events.AccessPointAddedMsg{AP: ap})
		},
		OnAccessPointRemoved: func(apPath godbus.ObjectPath) {
			p.Send(events.AccessPointRemovedMsg{ApPath: apPath})
		},
		OnNewConnection: func(path godbus.ObjectPath) {
			p.Send(events.NewConnectionMsg{ConnectionPath: path})
		},
		OnConnectionRemoved: func(path godbus.ObjectPath) {
			p.Send(events.ConnectionRemovedMsg{ConnectionPath: path})
		},
		OnActiveConnStateChanged: func(state, reason uint32) {
			if state == 2 {
				p.Send(events.UpdateActiveConnectionMsg{State: state, Reason: reason})
			}
			if state == 4 {
				p.Send(events.ClearActiveConnectionMsg{})
			}
		},
		OnError: func(err error) {
			p.Send(events.ShowError(err))
		},
	}
	go nm.ListenForSignals(client, cb)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
