package networkmanager

import godbus "github.com/godbus/dbus/v5"

// this is the shape of an Access Point.
// refer to the `.AccessPoint.Properties` in NM docs
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

// this is the shape of the saved connection profile
// as seen in `.Settings.Connection.GetSettings` (refer to NM docs)
type SavedConnection struct {
	// meta data
	ConnectionPath godbus.ObjectPath

	// connection block
	UUID        string
	AutoConnect bool // default is true

	// 802-11-wireless block
	SSID       string
	BSSIDs     []string
	Mode       string
	Band       string
	Hidden     bool
	MacAddress string

	// 802-11-wireless-security block
	KeyMgmt     string
	AuthAlg     string
	Proto       []string
	PskFlags    uint32
	WepKeyFlags uint32
	WepKeyType  uint32
	WpsMethod   uint32
}

// a saved connection that is currently visible
type NearbyConnection struct {
	Connection SavedConnection
	AP         *AccessPoint
}
