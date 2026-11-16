package networkmanager

// NearbyConnection is a saved connection that is currently in range
type NearbyConnection struct {
	Connection ConnectionSettings
	AP         *AccessPointProperties
}

type (
	// used for NMSettingSecretFlags bitmask
	SecretFlags uint32
	// used for NMWepKeyType
	WepKeyType uint32
	// used as the WpsMethod type in SavedConnection
	WpsMethod uint32
)
