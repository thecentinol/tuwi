package main

import (
	// "fmt"
	tea "charm.land/bubbletea/v2"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/ui"
	"log"
)

func main() {
	client, err := dbus.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	p := tea.NewProgram(ui.NewModel(client))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
