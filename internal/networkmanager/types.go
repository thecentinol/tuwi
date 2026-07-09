package networkmanager

import godbus "github.com/godbus/dbus/v5"

// this is the shape of an Access Point.
// refer to the `.AccessPoint.Properties` in NM docs
type AccessPoint struct {
	Flags    uint32
	WpaFlags uint32
	RsnFlags uint32

	SSID      []byte
	Frequency uint32
	BSSID     string // HwAddress

	Mode       uint32
	MaxBitrate uint32
	Bandwidth  uint32
	Strength   uint8
	LastSeen   int32

	DevicePath godbus.ObjectPath
	APPath     godbus.ObjectPath
}

// this is the shape of the saved connection profile
// as seen in `.Settings.Connection.GetSettings` (refer to NM docs)
type SavedConnection struct {
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
