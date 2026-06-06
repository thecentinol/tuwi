package dbus

const (
	baseServiceName = "org.freedesktop.NetworkManager"
	baseObjPath     = "/org/freedesktop/NetworkManager"
	deviceType      = baseServiceName + ".Device.DeviceType"
	deviceWireless  = baseServiceName + ".Device.Wireless"
	deviceWired     = baseServiceName + ".Device.Wired"

	accessPointFlags    = baseServiceName + ".AccessPoint.Flags"
	accessPointWpaFlags = baseServiceName + ".AccessPoint.WpaFlags"
	accessPointRsnFlags = baseServiceName + ".AccessPoint.RsnFlags"
	accessPointSsid     = baseServiceName + ".AccessPoint.Ssid"
	accessPointFreq     = baseServiceName + ".AccessPoint.Frequency"
	accessPointBssid    = baseServiceName + ".AccessPoint.HwAddress"
	accessPointMode     = baseServiceName + ".AccessPoint.Mode"
	accessPointStrength = baseServiceName + ".AccessPoint.Strength"
)
