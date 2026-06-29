package networkmanager

// Base Service + Object
const (
	BaseServiceName = "org.freedesktop.NetworkManager"
	BaseObjPath     = "/org/freedesktop/NetworkManager"
)

// Connection Manager - Base Object
const (
	// Methods
	CmReload                          = BaseServiceName + ".Reload"
	CmGetDevicesConnMan               = BaseServiceName + ".GetDevices"
	CmGetAllDevices                   = BaseServiceName + ".GetAllDevices"
	CmGetDeviceByIpIface              = BaseServiceName + ".GetDeviceByIpIface"
	CmActivateConnection              = BaseServiceName + ".ActivateConnection"
	CmAddAndActivateConnection2       = BaseServiceName + ".AddAndActivateConnection2"
	CmDeactivateConnection            = BaseServiceName + ".DeactivateConnection"
	CmSleep                           = BaseServiceName + ".Sleep"
	CmEnable                          = BaseServiceName + ".Enable"
	CmGetPermissions                  = BaseServiceName + ".GetPermissions"
	CmSetLogging                      = BaseServiceName + ".SetLogging"
	CmGetLogging                      = BaseServiceName + ".GetLogging"
	CmCheckConnectivity               = BaseServiceName + ".CheckConnectivity"
	CmState                           = BaseServiceName + ".state"
	CmCheckpointCreate                = BaseServiceName + ".CheckpointCreate"
	CmCheckpointDestroy               = BaseServiceName + ".CheckpointDestroy"
	CmCheckpointRollback              = BaseServiceName + ".CheckpointRollback"
	CmCheckpointAdjustRollbackTimeout = BaseServiceName + ".CheckpointAdjustRollbackTimeout"

	// Signals
	CmCheckPermissions = BaseServiceName + ".CheckPermissions"
	CmStateChanged     = BaseServiceName + ".StateChanged"
	CmDeviceAdded      = BaseServiceName + ".DeviceAdded"
	CmDeviceRemoved    = BaseServiceName + ".DeviceRemoved"

	// Properties
	CmDevices                    = BaseServiceName + ".Devices"
	CmAllDevices                 = BaseServiceName + ".AllDevices"
	CmCheckpoints                = BaseServiceName + ".Checkpoints"
	CmNetworkingEnabled          = BaseServiceName + ".NetworkingEnabled"
	CmWirelessEnabled            = BaseServiceName + ".WirelessEnabled"
	CmWirelessHardwareEnabled    = BaseServiceName + ".WirelessHardwareEnabled"
	CmWwanEnabled                = BaseServiceName + ".WwanEnabled"
	CmWwanHardwareEnabled        = BaseServiceName + ".WwanHardwareEnabled"
	CmWimaxEnabled               = BaseServiceName + ".WimaxEnabled"
	CmWimaxHardwareEnabled       = BaseServiceName + ".WimaxHardwareEnabled"
	CmRadioFlags                 = BaseServiceName + ".RadioFlags"
	CmActiveConnections          = BaseServiceName + ".ActiveConnections"
	CmPrimaryConnection          = BaseServiceName + ".PrimaryConnection"
	CmPrimaryConnectionType      = BaseServiceName + ".PrimaryConnectionType"
	CmMetered                    = BaseServiceName + ".Metered"
	CmActivatingConnection       = BaseServiceName + ".ActivatingConnection"
	CmStartup                    = BaseServiceName + ".Startup"
	CmVersion                    = BaseServiceName + ".Version"
	CmVersionInfo                = BaseServiceName + ".VersionInfo"
	CmCapabilities               = BaseServiceName + ".Capabilities"
	CmStateProperty              = BaseServiceName + ".State"
	CmConnectivity               = BaseServiceName + ".Connectivity"
	CmConnectivityCheckAvailable = BaseServiceName + ".ConnectivityCheckAvailable"
	CmConnectivityCheckEnabled   = BaseServiceName + ".ConnectivityCheckEnabled"
	CmConnectivityCheckUri       = BaseServiceName + ".ConnectivityCheckUri"
	CmGlobalDnsConfiguration     = BaseServiceName + ".GlobalDnsConfiguration"
)

// Connection Settings Profile Manager - Settings Object
const (
	SettingsBaseObjPath = BaseObjPath + "/Settings"

	// Methods
	CspmListConnections      = BaseServiceName + ".Settings.ListConnections"
	CspmGetConnectionByUuid  = BaseServiceName + ".Settings.GetConnectionByUuid"
	CspmAddConnection        = BaseServiceName + ".Settings.AddConnection"
	CspmAddConnectionUnsaved = BaseServiceName + ".Settings.AddConnectionUnsaved"
	CspmAddConnection2       = BaseServiceName + ".Settings.AddConnection2"
	CspmLoadConnections      = BaseServiceName + ".Settings.LoadConnections"
	CspmReloadConnections    = BaseServiceName + ".Settings.ReloadConnections"
	CspmSavedHostname        = BaseServiceName + ".Settings.SavedHostname"

	// Signals
	CspmNewConnection     = BaseServiceName + ".Settings.NewConnection"
	CspmConnectionRemoved = BaseServiceName + ".Settings.ConnectionRemoved"

	// Properties
	CspmConnections = BaseServiceName + ".Settings.Connections"
	CspmHostname    = BaseServiceName + ".Settings.Hostname"
	CspmCanModify   = BaseServiceName + ".Settings.CanModify"
	CspmVersionId   = BaseServiceName + ".Settings.VersionId"
)

// Connection Settings Profile - Settings Object
const (
	baseSettConnServiceName = BaseServiceName + ".Settings.Connection"

	// Methods
	CspUpdate        = baseSettConnServiceName + ".Update"
	CspUpdateUnsaved = baseSettConnServiceName + ".UpdateUnsaved"
	CspDelete        = baseSettConnServiceName + ".Delete"
	CspGetSettings   = baseSettConnServiceName + ".GetSettings"
	CspGetSecrets    = baseSettConnServiceName + ".GetSecrets"
	CspClearSecrets  = baseSettConnServiceName + ".ClearSecrets"
	CspSave          = baseSettConnServiceName + ".Save"
	CspUpdate2       = baseSettConnServiceName + ".Update2"

	// Signals
	CspUpdated = baseSettConnServiceName + ".Updated"
	CspRemoved = baseSettConnServiceName + ".Removed"

	// Properties
	CspUnsaved  = baseSettConnServiceName + ".Unsaved"
	CspFlags    = baseSettConnServiceName + ".Flags"
	CspFilename = baseSettConnServiceName + ".Filename"
)

// Device
const (
	baseDeviceServiceName = BaseServiceName + ".Device"

	// Methods
	DeviceReapply              = baseDeviceServiceName + ".Reapply"
	DeviceGetAppliedConnection = baseDeviceServiceName + ".GetAppliedConnection"
	DeviceDisconnect           = baseDeviceServiceName + ".Disconnect"
	DeviceDeleteDevice         = baseDeviceServiceName + ".Delete"

	// Signals
	DeviceStateChanged = baseDeviceServiceName + ".StateChanged"

	// Properties
	DeviceUdi                  = baseDeviceServiceName + ".Udi"
	DevicePath                 = baseDeviceServiceName + ".Path"
	DeviceInterface            = baseDeviceServiceName + ".Interface"
	DeviceIpInterface          = baseDeviceServiceName + ".IpInterface"
	DeviceDriver               = baseDeviceServiceName + ".Driver"
	DeviceDriverVersion        = baseDeviceServiceName + ".DriverVersion"
	DeviceFirmwareVersion      = baseDeviceServiceName + ".FirmwareVersion"
	DeviceCapabilities         = baseDeviceServiceName + ".Capabilities"
	DeviceIp4Address           = baseDeviceServiceName + ".Ip4Address"
	DeviceState                = baseDeviceServiceName + ".State"
	DeviceStateReason          = baseDeviceServiceName + ".StateReason"
	DeviceActiveConnection     = baseDeviceServiceName + ".ActiveConnection"
	DeviceIp4Config            = baseDeviceServiceName + ".Ip4Config"
	DeviceDhcp4Config          = baseDeviceServiceName + ".Dhcp4Config"
	DeviceIp6Config            = baseDeviceServiceName + ".Ip6Config"
	DeviceDhcp6Config          = baseDeviceServiceName + ".Dhcp6Config"
	DeviceManaged              = baseDeviceServiceName + ".Managed"
	DeviceAutoConnect          = baseDeviceServiceName + ".AutoConnect"
	DeviceFirmwareMissing      = baseDeviceServiceName + ".FirmwareMissing"
	DeviceNmPluginMissing      = baseDeviceServiceName + ".NmPluginMissing"
	DeviceType                 = baseDeviceServiceName + ".DeviceType"
	DeviceAvailableConnections = baseDeviceServiceName + ".AvailableConnections"
	DevicePhysicalPortId       = baseDeviceServiceName + ".PhysicalPortId"
	DeviceMtu                  = baseDeviceServiceName + ".Mtu"
	DeviceMetered              = baseDeviceServiceName + ".Metered"
	DeviceLldpNeighbors        = baseDeviceServiceName + ".LldpNeighbors"
	DeviceReal                 = baseDeviceServiceName + ".Real"
	DeviceIp4Connectivity      = baseDeviceServiceName + ".Ip4Connectivity"
	DeviceIp6Connectivity      = baseDeviceServiceName + ".Ip6Connectivity"
	DeviceInterfaceFlags       = baseDeviceServiceName + ".InterfaceFlags"
	DeviceHwAddress            = baseDeviceServiceName + ".HwAddress"
	DevicePorts                = baseDeviceServiceName + ".Ports"
)

// Device Wireless
const (
	baseWirelessServiceName = baseDeviceServiceName + ".Wireless"

	// Methods
	WirelessGetAllAccessPoints = baseWirelessServiceName + ".GetAllAccessPoints"
	WirelessRequestScan        = baseWirelessServiceName + ".RequestScan"

	// Signals
	WirelessAccessPointAdded   = baseWirelessServiceName + ".AccessPointAdded"
	WirelessAccessPointRemoved = baseWirelessServiceName + ".AccessPointRemoved"

	// Properties
	WirelessHwAddress            = baseWirelessServiceName + ".HwAddress"
	WirelessPermHwAddress        = baseWirelessServiceName + ".PermHwAddress"
	WirelessMode                 = baseWirelessServiceName + ".Mode"
	WirelessBitrate              = baseWirelessServiceName + ".Bitrate"
	WirelessAccessPoints         = baseWirelessServiceName + ".AccessPoints"
	WirelessActiveAccessPoint    = baseWirelessServiceName + ".ActiveAccessPoint"
	WirelessWirelessCapabilities = baseWirelessServiceName + ".WirelessCapabilities"
	WirelessLastScan             = baseWirelessServiceName + ".LastScan"
)

// Active Connection
const (
	baseActiveConnServiceName = BaseServiceName + ".Connection.Active"

	// Signals
	AcStateChanged = baseActiveConnServiceName + ".StateChanged"

	// Properties
	ActConConnection     = baseActiveConnServiceName + ".Connection"
	ActConSpecificObject = baseActiveConnServiceName + ".SpecificObject"
	ActConId             = baseActiveConnServiceName + ".Id"
	ActConUuid           = baseActiveConnServiceName + ".Uuid"
	ActConType           = baseActiveConnServiceName + ".Type"
	ActConDevices        = baseActiveConnServiceName + ".Devices"
	ActConState          = baseActiveConnServiceName + ".State"
	ActConStateFlags     = baseActiveConnServiceName + ".StateFlags"
	ActConDefault        = baseActiveConnServiceName + ".Default"
	ActConIp4Config      = baseActiveConnServiceName + ".Ip4Config"
	ActConDhcp4Config    = baseActiveConnServiceName + ".Dhcp4Config"
	ActConDefault6       = baseActiveConnServiceName + ".Default6"
	ActConIp6Config      = baseActiveConnServiceName + ".Ip6Config"
	ActConDhcp6Config    = baseActiveConnServiceName + ".Dhcp6Config"
	ActConVpn            = baseActiveConnServiceName + ".Vpn"
	ActConController     = baseActiveConnServiceName + ".Controller"
	ActConMaster         = baseActiveConnServiceName + ".Master"
)

// Access Point Flags, WpaFlags, RsnFlags
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

// Wi-Fi Access Point
const (
	baseApServiceName = BaseServiceName + ".AccessPoint"

	// Properties
	ApFlags      = baseApServiceName + ".Flags"
	ApWpaFlags   = baseApServiceName + ".WpaFlags"
	ApRsnFlags   = baseApServiceName + ".RsnFlags"
	ApSsid       = baseApServiceName + ".Ssid"
	ApFrequency  = baseApServiceName + ".Frequency"
	ApHwAddress  = baseApServiceName + ".HwAddress"
	ApMode       = baseApServiceName + ".Mode"
	ApMaxBitrate = baseApServiceName + ".MaxBitrate"
	ApBandwidth  = baseApServiceName + ".Bandwidth"
	ApStrength   = baseApServiceName + ".Strength"
	ApLastSeen   = baseApServiceName + ".LastSeen"
)
