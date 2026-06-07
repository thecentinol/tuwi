package dbus

const (
	baseServiceName = "org.freedesktop.NetworkManager"
	baseObjPath     = "/org/freedesktop/NetworkManager"

	deviceType = baseServiceName + ".Device.DeviceType"

	// wireless
	deviceWireless    = baseServiceName + ".Device.Wireless"
	activeAccessPoint = baseServiceName + ".Device.Wireless.ActiveAccessPoint"

	// wired
	deviceWired = baseServiceName + ".Device.Wired"

	// access point
	accessPointFlags    = baseServiceName + ".AccessPoint.Flags"
	accessPointWpaFlags = baseServiceName + ".AccessPoint.WpaFlags"
	accessPointRsnFlags = baseServiceName + ".AccessPoint.RsnFlags"
	accessPointSsid     = baseServiceName + ".AccessPoint.Ssid"
	accessPointFreq     = baseServiceName + ".AccessPoint.Frequency"
	accessPointBssid    = baseServiceName + ".AccessPoint.HwAddress"
	accessPointMode     = baseServiceName + ".AccessPoint.Mode"
	accessPointStrength = baseServiceName + ".AccessPoint.Strength"
)
