package events

import (
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// Device Wireless
type AccessPointAddedMsg struct {
	AP nm.AccessPoint
}
type AccessPointRemovedMsg struct {
	ApPath string
}

// Connection Settings Profile Manager - Settings Object
type NewConnectionMsg struct {
	ConnectionPath string
}
type ConnectionRemovedMsg struct {
	ConnectionPath string
}
