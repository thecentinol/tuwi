package main

import (
	tea "charm.land/bubbletea/v2"
	"log"

	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/events"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/ui"
)

func main() {
	client, err := dbus.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	p := tea.NewProgram(ui.NewModel(client))

	cb := nm.SignalCallbacks{
		OnAccessPointAdded: func(ap nm.AccessPoint) {
			p.Send(events.AccessPointAddedMsg{AP: ap})
		},
		OnAccessPointRemoved: func(apPath string) {
			p.Send(events.AccessPointRemovedMsg{ApPath: apPath})
		},
		OnNewConnection: func(path string) {
			p.Send(events.NewConnectionMsg{ConnectionPath: path})
		},
		OnConnectionRemoved: func(path string) {
			p.Send(events.ConnectionRemovedMsg{ConnectionPath: path})
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
