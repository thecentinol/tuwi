package networkmanager

// Base Service + Object
const (
	BaseServiceName = "org.freedesktop.NetworkManager"
	BaseObjPath     = "/org/freedesktop/NetworkManager"
)

// Connection Manager - Base Object
const (
	// Methods
	Reload                          = BaseServiceName + ".Reload"
	GetDevices                      = BaseServiceName + ".GetDevices"
	GetAllDevices                   = BaseServiceName + ".GetAllDevices"
	GetDeviceByIpIface              = BaseServiceName + ".GetDeviceByIpIface"
	ActivateConnection              = BaseServiceName + ".ActivateConnection"
	AddAndActivateConnection2       = BaseServiceName + ".AddAndActivateConnection2"
	DeactivateConnection            = BaseServiceName + ".DeactivateConnection"
	Sleep                           = BaseServiceName + ".Sleep"
	Enable                          = BaseServiceName + ".Enable"
	GetPermissions                  = BaseServiceName + ".GetPermissions"
	SetLogging                      = BaseServiceName + ".SetLogging"
	GetLogging                      = BaseServiceName + ".GetLogging"
	CheckConnectivity               = BaseServiceName + ".CheckConnectivity"
	StateMethod                     = BaseServiceName + ".state"
	CheckpointCreate                = BaseServiceName + ".CheckpointCreate"
	CheckpointDestroy               = BaseServiceName + ".CheckpointDestroy"
	CheckpointRollback              = BaseServiceName + ".CheckpointRollback"
	CheckpointAdjustRollbackTimeout = BaseServiceName + ".CheckpointAdjustRollbackTimeout"

	// Signals
	CheckPermissions = BaseServiceName + ".CheckPermissions"
	StateChanged     = BaseServiceName + ".StateChanged"
	DeviceAdded      = BaseServiceName + ".DeviceAdded"
	DeviceRemoved    = BaseServiceName + ".DeviceRemoved"

	// Properties
	Devices                    = BaseServiceName + ".Devices"
	AllDevices                 = BaseServiceName + ".AllDevices"
	Checkpoints                = BaseServiceName + ".Checkpoints"
	NetworkingEnabled          = BaseServiceName + ".NetworkingEnabled"
	WirelessEnabled            = BaseServiceName + ".WirelessEnabled"
	WirelessHardwareEnabled    = BaseServiceName + ".WirelessHardwareEnabled"
	WwanEnabled                = BaseServiceName + ".WwanEnabled"
	WwanHardwareEnabled        = BaseServiceName + ".WwanHardwareEnabled"
	WimaxEnabled               = BaseServiceName + ".WimaxEnabled"
	WimaxHardwareEnabled       = BaseServiceName + ".WimaxHardwareEnabled"
	RadioFlags                 = BaseServiceName + ".RadioFlags"
	ActiveConnections          = BaseServiceName + ".ActiveConnections"
	PrimaryConnection          = BaseServiceName + ".PrimaryConnection"
	PrimaryConnectionType      = BaseServiceName + ".PrimaryConnectionType"
	Metered                    = BaseServiceName + ".Metered"
	ActivatingConnection       = BaseServiceName + ".ActivatingConnection"
	Startup                    = BaseServiceName + ".Startup"
	Version                    = BaseServiceName + ".Version"
	VersionInfo                = BaseServiceName + ".VersionInfo"
	Capabilities               = BaseServiceName + ".Capabilities"
	StateProperty              = BaseServiceName + ".State"
	Connectivity               = BaseServiceName + ".Connectivity"
	ConnectivityCheckAvailable = BaseServiceName + ".ConnectivityCheckAvailable"
	ConnectivityCheckEnabled   = BaseServiceName + ".ConnectivityCheckEnabled"
	ConnectivityCheckUri       = BaseServiceName + ".ConnectivityCheckUri"
	GlobalDnsConfiguration     = BaseServiceName + ".GlobalDnsConfiguration"
)

// Connection Settings Profile Manager - Settings Object
const (
	SettingsBaseObjPath = BaseObjPath + "/Settings"

	// Methods
	ListConnections      = BaseServiceName + ".Settings.ListConnections"
	GetConnectionByUuid  = BaseServiceName + ".Settings.GetConnectionByUuid"
	AddConnection        = BaseServiceName + ".Settings.AddConnection"
	AddConnectionUnsaved = BaseServiceName + ".Settings.AddConnectionUnsaved"
	AddConnection2       = BaseServiceName + ".Settings.AddConnection2"
	LoadConnections      = BaseServiceName + ".Settings.LoadConnections"
	ReloadConnections    = BaseServiceName + ".Settings.ReloadConnections"
	SavedHostname        = BaseServiceName + ".Settings.SavedHostname"

	// Signals
	NewConnection     = BaseServiceName + ".Settings.NewConnection"
	ConnectionRemoved = BaseServiceName + ".Settings.ConnectionRemoved"

	// Properties
	Connections = BaseServiceName + ".Settings.Connections"
	Hostname    = BaseServiceName + ".Settings.Hostname"
	CanModify   = BaseServiceName + ".Settings.CanModify"
	VersionId   = BaseServiceName + ".Settings.VersionId"
)

// Connection Settings Profile - Settings Object
const (
	baseSettConnServiceName = BaseServiceName + ".Settings.Connection"

	// Methods
	Update                   = baseSettConnServiceName + ".Update"
	UpdateUnsaved            = baseSettConnServiceName + ".UpdateUnsaved"
	DeleteSettingsConnection = baseSettConnServiceName + ".Delete"
	GetSettings              = baseSettConnServiceName + ".GetSettings"
	GetSecrets               = baseSettConnServiceName + ".GetSecrets"
	ClearSecrets             = baseSettConnServiceName + ".ClearSecrets"
	Save                     = baseSettConnServiceName + ".Save"
	Update2                  = baseSettConnServiceName + ".Update2"

	// Signals
	Updated = baseSettConnServiceName + ".Updated"
	Removed = baseSettConnServiceName + ".Removed"

	// Properties
	Unsaved  = baseSettConnServiceName + ".Unsaved"
	Flags    = baseSettConnServiceName + ".Flags"
	Filename = baseSettConnServiceName + ".Filename"
)

// Device
const (
	baseDeviceServiceName     = BaseServiceName + ".Device"
	deviceWirelessServiceName = baseDeviceServiceName + ".Wireless"

	// Methods
	Reapply              = baseDeviceServiceName + ".Reapply"
	GetAppliedConnection = baseDeviceServiceName + ".GetAppliedConnection"
	Disconnect           = baseDeviceServiceName + ".Disconnect"
	DeleteDevice         = baseDeviceServiceName + ".Delete"

	// Signals
	StateChangedDevice = baseDeviceServiceName + ".StateChanged"

	// Properties
	Udi                  = baseDeviceServiceName + ".Udi"
	Path                 = baseDeviceServiceName + ".Path"
	Interface            = baseDeviceServiceName + ".Interface"
	IpInterface          = baseDeviceServiceName + ".IpInterface"
	Driver               = baseDeviceServiceName + ".Driver"
	DriverVersion        = baseDeviceServiceName + ".DriverVersion"
	FirmwareVersion      = baseDeviceServiceName + ".FirmwareVersion"
	CapabilitiesDevice   = baseDeviceServiceName + ".Capabilities"
	Ip4Address           = baseDeviceServiceName + ".Ip4Address"
	StateDevice          = baseDeviceServiceName + ".State"
	StateReason          = baseDeviceServiceName + ".StateReason"
	ActiveConnection     = baseDeviceServiceName + ".ActiveConnection"
	Ip4Config            = baseDeviceServiceName + ".Ip4Config"
	Dhcp4Config          = baseDeviceServiceName + ".Dhcp4Config"
	Ip6Config            = baseDeviceServiceName + ".Ip6Config"
	Dhcp6Config          = baseDeviceServiceName + ".Dhcp6Config"
	Managed              = baseDeviceServiceName + ".Managed"
	AutoConnect          = baseDeviceServiceName + ".AutoConnect"
	FirmwareMissing      = baseDeviceServiceName + ".FirmwareMissing"
	NmPluginMissing      = baseDeviceServiceName + ".NmPluginMissing"
	DeviceType           = baseDeviceServiceName + ".DeviceType"
	AvailableConnections = baseDeviceServiceName + ".AvailableConnections"
	PhysicalPortId       = baseDeviceServiceName + ".PhysicalPortId"
	Mtu                  = baseDeviceServiceName + ".Mtu"
	MeteredDevice        = baseDeviceServiceName + ".Metered"
	LldpNeighbors        = baseDeviceServiceName + ".LldpNeighbors"
	Real                 = baseDeviceServiceName + ".Real"
	Ip4Connectivity      = baseDeviceServiceName + ".Ip4Connectivity"
	Ip6Connectivity      = baseDeviceServiceName + ".Ip6Connectivity"
	InterfaceFlags       = baseDeviceServiceName + ".InterfaceFlags"
	HwAddress            = baseDeviceServiceName + ".HwAddress"
	Ports                = baseDeviceServiceName + ".Ports"
)

// Device Wireless
const (
	baseWirelessServiceName = baseDeviceServiceName + ".Wireless"

	// Methods
	GetAllAccessPoints = baseWirelessServiceName + ".GetAllAccessPoints"
	RequestScan        = baseWirelessServiceName + ".RequestScan"

	// Signals
	AccessPointAdded   = baseWirelessServiceName + ".AccessPointAdded"
	AccessPointRemoved = baseWirelessServiceName + ".AccessPointRemoved"

	// Properties
	HwAddressWireless    = baseWirelessServiceName + ".HwAddress"
	PermHwAddress        = baseWirelessServiceName + ".PermHwAddress"
	Mode                 = baseWirelessServiceName + ".Mode"
	Bitrate              = baseWirelessServiceName + ".Bitrate"
	AccessPoints         = baseWirelessServiceName + ".AccessPoints"
	ActiveAccessPoint    = baseWirelessServiceName + ".ActiveAccessPoint"
	WirelessCapabilities = baseWirelessServiceName + ".WirelessCapabilities"
	LastScan             = baseWirelessServiceName + ".LastScan"
)

// Active Connection
const (
	baseActiveConnServiceName = BaseServiceName + ".Connection.Active"

	// Signals
	StateChangedActiveConn = baseActiveConnServiceName + ".StateChanged"

	// Properties
	Connection              = baseActiveConnServiceName + ".Connection"
	SpecificObject          = baseActiveConnServiceName + ".SpecificObject"
	Id                      = baseActiveConnServiceName + ".Id"
	Uuid                    = baseActiveConnServiceName + ".Uuid"
	TypeActiveConnection    = baseActiveConnServiceName + ".Type"
	DevicesActiveConnection = baseActiveConnServiceName + ".Devices"
	StateActiveConnection   = baseActiveConnServiceName + ".State"
	StateFlagsActiveConn    = baseActiveConnServiceName + ".StateFlags"
	DefaultActiveConnection = baseActiveConnServiceName + ".Default"
	Ip4ConfigActiveConn     = baseActiveConnServiceName + ".Ip4Config"
	Dhcp4ConfigActiveConn   = baseActiveConnServiceName + ".Dhcp4Config"
	Default6                = baseActiveConnServiceName + ".Default6"
	Ip6ConfigActiveConn     = baseActiveConnServiceName + ".Ip6Config"
	Dhcp6ConfigActiveConn   = baseActiveConnServiceName + ".Dhcp6Config"
	Vpn                     = baseActiveConnServiceName + ".Vpn"
	Controller              = baseActiveConnServiceName + ".Controller"
	Master                  = baseActiveConnServiceName + ".Master"
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
	FlagsAp     = baseApServiceName + ".Flags"
	WpaFlags    = baseApServiceName + ".WpaFlags"
	RsnFlags    = baseApServiceName + ".RsnFlags"
	Ssid        = baseApServiceName + ".Ssid"
	Frequency   = baseApServiceName + ".Frequency"
	HwAddressAp = baseApServiceName + ".HwAddress"
	ModeAp      = baseApServiceName + ".Mode"
	MaxBitrate  = baseApServiceName + ".MaxBitrate"
	Bandwidth   = baseApServiceName + ".Bandwidth"
	Strength    = baseApServiceName + ".Strength"
	LastSeen    = baseApServiceName + ".LastSeen"
)
