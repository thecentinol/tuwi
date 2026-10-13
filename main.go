package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/config"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	"github.com/thecentinol/tuwi/internal/keybindings"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/ui/root"
)

func main() {
	client, err := dbus.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// load config
	cfg, errCfg := config.Load()
	if !errors.Is(errCfg, os.ErrNotExist) && !errors.Is(errCfg, os.ErrPermission) && errCfg != nil {
		log.Fatal(err)
	}

	keys := keybindings.DefaultKeybindings()
	model := root.NewModel(client, keys, *cfg)

	if errors.Is(errCfg, os.ErrPermission) {
		model = model.SetStartupError(
			fmt.Errorf("permission error when trying to read config file, falling back to default config"),
		)
	}

	p := tea.NewProgram(model)

	cb := nm.SignalCallbacks{
		OnAccessPointAdded: func(ap nm.AccessPoint) {
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
		log.Fatal(err)
	}
}
