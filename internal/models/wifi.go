package models

import (
	godbus "github.com/godbus/dbus/v5"
)

type AccessPoint struct {
	Hidden   bool
	SSID     string
	BSSID    string // AKA HwAddress
	Strength uint8
	Secured  bool // derived from privacy flag
	HasWps   bool

	DevicePath godbus.ObjectPath
	APPath     godbus.ObjectPath
}
