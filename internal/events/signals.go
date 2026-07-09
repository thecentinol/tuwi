package events

import (
	godbus "github.com/godbus/dbus/v5"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// Device Wireless
type AccessPointAddedMsg struct {
	AP nm.AccessPoint
}
type AccessPointRemovedMsg struct {
	ApPath godbus.ObjectPath
}

// Connection Settings Profile Manager - Settings Object
type NewConnectionMsg struct {
	ConnectionPath godbus.ObjectPath
}
type ConnectionRemovedMsg struct {
	ConnectionPath godbus.ObjectPath
}
