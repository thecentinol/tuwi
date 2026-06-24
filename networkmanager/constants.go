package networkmanager

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
	AddAndConnect     = BaseServiceName + ".AddAndActivateConnection2"

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

// //////////////////////////////////////////
// Access Point Flags, WpaFlags, RsnFlags //
// //////////////////////////////////////////
const (
	// AP Flags
	NmApFlagsNone    = 0x00000000
	NmApFlagsPrivacy = 0x00000001
	NmApFlagsWps     = 0x00000002
	NmApFlagsWpsPbc  = 0x00000004
	NmApFlagsWpsPin  = 0x00000008

	// AP Security Flags
	NmSecMgmtPsk       = 0x00000100
	NmSecMgmt8021      = 0x00000200
	NmSecMgmtSae       = 0x00000400
	NmSecMgmtOwe       = 0x00000800
	NmSecMgmtOweTm     = 0x00001000
	NmSecMgmtSuiteB192 = 0x00002000 // WPA3 Enterprise Suite-B
)
