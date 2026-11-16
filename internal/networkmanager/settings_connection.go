package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

const (
	settingsConnectionIface = settingsIface + ".Connection"

	// Methods
	settingsConnectionDelete      = settingsConnectionIface + ".Delete"
	settingsConnectionGetSettings = settingsConnectionIface + ".GetSettings"

	// NMSettingSecretFlags
	SettingsConnectionSettingSecretFlagNone        SecretFlags = 0
	SettingsConnectionSettingSecretFlagAgentOwned  SecretFlags = 1 << 0
	SettingsConnectionSettingSecretFlagNotSaved    SecretFlags = 1 << 1
	SettingsConnectionSettingSecretFlagNotRequired SecretFlags = 1 << 2
	SettingsConnectionSettingSecretFlagsKnown                  = SettingsConnectionSettingSecretFlagAgentOwned |
		SettingsConnectionSettingSecretFlagNotSaved | SettingsConnectionSettingSecretFlagNotRequired

	// NMWepKeyType
	settingsConnectionWepKeyTypeUnknown    WepKeyType = 0
	settingsConnectionWepKeyTypeKey        WepKeyType = 1
	settingsConnectionWepKeyTypePassphrase WepKeyType = 2

	// wps-method flags
	settingsConnectionWpsMethodDefault    WpsMethod = 0
	settingsConnectionWpsMethodDisabled   WpsMethod = 1 << 0
	settingsConnectionWpsMethodAuto       WpsMethod = 1 << 1
	settingsConnectionWpsMethodPbc        WpsMethod = 1 << 2
	settingsConnectionWpsMethodPin        WpsMethod = 1 << 3
	settingsConnectionWpsMethodPinDisplay WpsMethod = 1 << 4

	// the top level block's of settings map returned from ...Settings.Connection.GetSettings().
	// the map type := map[string]map[string]godbus.Variant
	ConnectionSettingsConnectionBlock       = "connection"               // the connection block
	ConnectionSettingsWirelessBlock         = "802-11-wireless"          // the 802-11-wirless block
	ConnectionSettingsWirelessSecurityBlock = "802-11-wireless-security" // the 802-11-wirless-security block
	ConnectionSettingsIpv4Block             = "ipv4"                     // the IPv4 block
	ConnectionSettingsIpv6Block             = "ipv6"                     // the IPv6 block
)

var SettingsConnection settingsConnection

type settingsConnection struct{}

// return type from ...Settings.Connection.GetSettings().
// Does not include all settings, just the ones that are needed/wanted.
type ConnectionSettings struct {
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
	WepKeyType  WepKeyType
	WpsMethod   WpsMethod
}

func (settingsConnection) Delete(c *dbus.Client, conPath godbus.ObjectPath) error {
	obj := c.Conn.Object(baseIface, conPath)
	call := obj.Call(settingsConnectionDelete, 0)

	if call.Err != nil {
		return fmt.Errorf("DeleteConnection: %w", call.Err)
	}

	return nil
}

func (settingsConnection) GetSettings(c *dbus.Client, path godbus.ObjectPath) (*ConnectionSettings, error) {
	// as per NM docs this is the response type from GetSettings:
	var settings map[string]map[string]godbus.Variant
	var uuid string
	autoConnect := true
	ssid := "unknown"
	var bssids []string
	var mode string
	var band string
	var hidden bool
	var macAddress string
	var keyMgmt string
	var authAlg string
	var proto []string
	var pskFlags uint32
	var wepKeyFlags uint32
	var wepKeyType uint32
	var wpsMethod uint32

	obj := c.Conn.Object(baseIface, path)
	err := obj.Call(settingsConnectionGetSettings, 0).Store(&settings)
	if err != nil {
		return nil, fmt.Errorf("GetSettings: store: %w", err)
	}

	// filter out non-WiFi connections
	connBlock, ok := settings["connection"]
	if !ok || connBlock["type"].Value().(string) != ConnectionSettingsWirelessBlock {
		return nil, nil
	}

	rawUuid, ok := connBlock["uuid"]
	if ok {
		uuid = rawUuid.Value().(string)
	}

	autoConn, ok := connBlock["autoconnect"]
	if ok {
		autoConnect = autoConn.Value().(bool)
	}

	wirelessBlock, ok := settings["802-11-wireless"]
	if !ok {
		return nil, fmt.Errorf("GetSettings: 802-11-wireless block not found")
	}

	rawSsid, ok := wirelessBlock["ssid"]
	if ok {
		if ssidBytes, ok := rawSsid.Value().([]byte); ok {
			ssid = string(ssidBytes)
		}
	}

	seenBssids, ok := wirelessBlock["seen-bssids"]
	if ok {
		bssids = seenBssids.Value().([]string)
	}

	rawMode, ok := wirelessBlock["mode"]
	if ok {
		mode = rawMode.Value().(string)
	}

	rawBand, ok := wirelessBlock["band"]
	if ok {
		band = rawBand.Value().(string)
	}

	rawHidden, ok := wirelessBlock["hidden"]
	if ok {
		hidden = rawHidden.Value().(bool)
	}

	rawMac, ok := wirelessBlock["mac-address"]
	if ok {
		macAddress = rawMac.Value().(string)
	}

	securityBlock, ok := settings[ConnectionSettingsWirelessSecurityBlock]
	if ok {
		rawKeyMgmt, ok := securityBlock["key-mgmt"]
		if ok {
			keyMgmt = rawKeyMgmt.Value().(string)
		}

		rawAuthAlg, ok := securityBlock["auth-alg"]
		if ok {
			authAlg = rawAuthAlg.Value().(string)
		}

		rawProto, ok := securityBlock["proto"]
		if ok {
			proto = rawProto.Value().([]string)
		}

		rawPsk, ok := securityBlock["psk-flags"]
		if ok {
			pskFlags = rawPsk.Value().(uint32)
		}

		rawWepKeyFlags, ok := securityBlock["wep-key-flags"]
		if ok {
			wepKeyFlags = rawWepKeyFlags.Value().(uint32)
		}

		rawWepKeyType, ok := securityBlock["wep-key-type"]
		if ok {
			wepKeyType = rawWepKeyType.Value().(uint32)
		}

		rawWpsMeth, ok := securityBlock["wps-method"]
		if ok {
			wpsMethod = rawWpsMeth.Value().(uint32)
		}
	}

	sc := ConnectionSettings{
		ConnectionPath: path,
		UUID:           uuid,
		AutoConnect:    autoConnect,
		SSID:           ssid,
		BSSIDs:         bssids,
		Mode:           mode,
		Band:           band,
		Hidden:         hidden,
		MacAddress:     macAddress,
		KeyMgmt:        keyMgmt,
		AuthAlg:        authAlg,
		Proto:          proto,
		PskFlags:       pskFlags,
		WepKeyFlags:    wepKeyFlags,
		WepKeyType:     WepKeyType(wepKeyType),
		WpsMethod:      WpsMethod(wpsMethod),
	}

	return &sc, nil
}

// Has checks if a SecretFlags has the given flag
func (f SecretFlags) Has(flag SecretFlags) bool { return f&flag != 0 }

// String returns the string representation of a NMWepKeyType
func (t WepKeyType) String() string {
	switch t {
	case settingsConnectionWepKeyTypeUnknown:
		return "Unknown"
	case settingsConnectionWepKeyTypeKey:
		return "Key"
	case settingsConnectionWepKeyTypePassphrase:
		return "Passphrase"
	}

	return fmt.Sprintf("Unknown (%d)", uint32(t))
}

// String returns the string representation of a WpsMethod
func (m WpsMethod) String() string {
	switch m {
	case settingsConnectionWpsMethodDefault:
		return "Default"
	case settingsConnectionWpsMethodDisabled:
		return "Disabled"
	case settingsConnectionWpsMethodAuto:
		return "Auto"
	case settingsConnectionWpsMethodPbc:
		return "Push Button Connect"
	case settingsConnectionWpsMethodPin:
		return "PIN"
	case settingsConnectionWpsMethodPinDisplay:
		return "PIN (display)"
	}

	return fmt.Sprintf("0x%x", uint32(m))
}
