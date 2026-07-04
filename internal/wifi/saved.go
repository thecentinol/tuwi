package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// get a slice of the saved connection profile's
func GetSavedNetworks(c *dbus.Client, available []nm.AccessPoint) ([]nm.AccessPoint, error) {
	settingsObj := c.Conn.Object(nm.BaseServiceName, nm.SettingsBaseObjPath)

	var connPaths []godbus.ObjectPath
	err := settingsObj.Call(nm.CspmListConnections, 0).Store(&connPaths)
	if err != nil {
		return nil, fmt.Errorf("GetSavedNetworks: list connections: %w", err)
	}

	var savedNetworks []nm.AccessPoint

	for _, path := range connPaths {
		ap, err := GetApFromSettings(c, path, available)
		if err != nil {
			return nil, fmt.Errorf("GetSavedNetworks: %w", err)
		}

		if ap != nil {
			savedNetworks = append(savedNetworks, *ap)
		}
	}

	return savedNetworks, nil
}

// build an AccessPoint from a saved connection profile's settings.
// `available` is the slice of available networks that we already
// have so that we can cross-reference/match the new connection with an existing one,
// this is done so we can get info on the network that's not available via it's
// settings - such as Strength.
func GetApFromSettings(c *dbus.Client, path godbus.ObjectPath, available []nm.AccessPoint) (*nm.AccessPoint, error) {
	// as per NM docs this is the response type from GetSettings.
	var settings map[string]map[string]godbus.Variant

	var ssid string = "unknown"
	var bssid string
	var conUUID string = ""
	var secType string = "open"
	var strength uint8 = 0
	var activeAPPath godbus.ObjectPath
	var devicePath godbus.ObjectPath

	visibleAPs := make(map[string]nm.AccessPoint)
	for _, ap := range available {
		if ap.BSSID != "" {
			visibleAPs[ap.BSSID] = ap
		}
	}

	obj := c.Conn.Object(nm.BaseServiceName, godbus.ObjectPath(path))
	err := obj.Call(nm.CspGetSettings, 0).Store(&settings)
	if err != nil {
		return nil, fmt.Errorf("GetApFromSettings: GetSettings: %w", err)
	}

	// filter out non-WiFi connections
	connBlock, ok := settings["connection"]
	if !ok || connBlock["type"].Value().(string) != "802-11-wireless" {
		return nil, nil
	}

	uuid, ok := connBlock["uuid"]
	if ok {
		conUUID = uuid.Value().(string)
	}

	if wirelessBlock, ok := settings["802-11-wireless"]; ok {
		// get SSID
		if ssidVal, ok := wirelessBlock["ssid"]; ok {
			if ssidBytes, ok := ssidVal.Value().([]byte); ok {
				ssid = string(ssidBytes)
			}
		}

		// get BSSID
		if bssidVal, ok := wirelessBlock["seen-bssids"]; ok {
			if bssidSlice, ok := bssidVal.Value().([]string); ok {
				for _, seenBssid := range bssidSlice {
					if ap, found := visibleAPs[seenBssid]; found {
						strength = ap.Strength
						activeAPPath = ap.APPath
						devicePath = ap.DevicePath
						bssid = ap.BSSID

						break
					}
				}
			}
		}
	}

	_, secured := settings["802-11-wireless-security"]
	if secured {
		if keyMgmtVal, ok := settings["802-11-wireless-security"]["key-mgmt"]; ok {
			secType = keyMgmtVal.Value().(string)
		}
	}

	ap := nm.AccessPoint{
		SSID:           ssid,
		BSSID:          bssid,
		ConnectionUUID: conUUID,
		SecurityType:   secType,

		ConnectionPath: path,
		DevicePath:     devicePath,
		APPath:         activeAPPath,

		Strength: strength,

		IsSaved: true,
		Secured: secured,
		Hidden:  false,
	}

	return &ap, nil
}

// get the wifi network that's currently connected to
func GetActiveNetwork(c *dbus.Client) (*nm.AccessPoint, error) {
	var activeNetwork nm.AccessPoint

	devicePath, err := nm.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}
	device := c.Conn.Object(nm.BaseServiceName, devicePath)
	active, err := device.GetProperty(nm.WirelessActiveAccessPoint)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActiveAccessPoint property: %w", err)
	}
	activePath := active.Value().(godbus.ObjectPath)
	apObject := c.Conn.Object(nm.BaseServiceName, activePath)

	ssid, err := apObject.GetProperty(nm.ApSsid)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActivePointSsid property: %w", err)
	}

	strength, err := apObject.GetProperty(nm.ApStrength)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: AccessPointStrength property: %w", err)
	}

	flags, err := apObject.GetProperty(nm.ApFlags)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: AccessPointFlags property: %w", err)
	}

	ssidBytes := string(ssid.Value().([]byte))
	flagsVal := flags.Value().(uint32)
	activeNetwork = nm.AccessPoint{
		SSID:     ssidBytes,
		Strength: strength.Value().(uint8),
		Secured:  flagsVal&0x00000001 != 0,
		HasWps:   flagsVal&0x00000002 != 0,
	}

	return &activeNetwork, nil
}
