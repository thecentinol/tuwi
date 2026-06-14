package dbus

const (
	BaseServiceName = "org.freedesktop.NetworkManager"
	BaseObjPath     = "/org/freedesktop/NetworkManager"

	DeviceType = BaseServiceName + ".Device.DeviceType"

	// settings
	SettingsBaseObjPath     = BaseObjPath + "/Settings"
	SettingsBaseServiceName = BaseServiceName + ".Settings"
	ListConnections         = SettingsBaseServiceName + ".ListConnections"

	// settings connection
	SettConnection = SettingsBaseServiceName + ".Connection"
	GetSettings    = SettConnection + ".GetSettings"

	// wireless
	DeviceWireless    = BaseServiceName + ".Device.Wireless"
	ActiveAccessPoint = BaseServiceName + ".Device.Wireless.ActiveAccessPoint"

	// wired
	DeviceWired = BaseServiceName + ".Device.Wired"

	// access point
	AccessPointFlags    = BaseServiceName + ".AccessPoint.Flags"
	AccessPointWpaFlags = BaseServiceName + ".AccessPoint.WpaFlags"
	AccessPointRsnFlags = BaseServiceName + ".AccessPoint.RsnFlags"
	AccessPointSsid     = BaseServiceName + ".AccessPoint.Ssid"
	AccessPointFreq     = BaseServiceName + ".AccessPoint.Frequency"
	AccessPointBssid    = BaseServiceName + ".AccessPoint.HwAddress"
	AccessPointMode     = BaseServiceName + ".AccessPoint.Mode"
	AccessPointStrength = BaseServiceName + ".AccessPoint.Strength"
)
