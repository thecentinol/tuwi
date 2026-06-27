package networkmanager

import godbus "github.com/godbus/dbus/v5"

type AccessPoint struct {
	SSID           string
	BSSID          string // AKA HwAddress
	ConnectionUUID string
	SecurityType   string // wpa-psk, wpa-eap, none

	ConnectionPath godbus.ObjectPath
	DevicePath     godbus.ObjectPath
	APPath         godbus.ObjectPath

	Strength uint8

	IsSaved bool
	Secured bool
	Hidden  bool
	HasWps  bool
}
