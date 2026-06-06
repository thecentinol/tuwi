package main

import (
	"fmt"
	"log"
	"tuwi/internal/dbus"
)

func main() {
	client, err := dbus.NewClient()
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	defer client.Close()

	devices, err := dbus.GetDevices(client)
	if err != nil {
		log.Fatalf("Error fetching devices: %v", err)
	}

	fmt.Printf("\nFound %d networks:\n", len(devices))
	for _, devs := range devices {
		fmt.Printf("Devices: %v\n", devs)
	}

	aps, err := dbus.GetAccessPoints(client)
	if err != nil {
		log.Fatalf("Error fetch access points: &v", err)
	}
	for _, v := range aps {
		fmt.Printf("Access Point: %v", v)
	}
}
