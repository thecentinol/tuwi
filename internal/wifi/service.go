package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// this gets the actual APs + details (ssid, strength, etc)
// for available WiFi networks
func GetAvailableNetworks(c *dbus.Client) ([]nm.AccessPoint, error) {
	wifiDevicePath, err := nm.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetWifiDevice: %w", err)
	}

	if err := nm.RequestScan(c, wifiDevicePath); err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: %w", err)
	}

	accessPoints, err := nm.GetAccessPoints(c, wifiDevicePath)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetAccessPoints: %w", err)
	}

	var networks []nm.AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.Conn.Object(nm.BaseServiceName, accessPointPath)

		ssidBytes, err := ap.GetProperty(nm.ApSsid)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointSsid property: %w", err)
		}
		rawSsid := ssidBytes.Value().([]byte)

		hidden, ssid := nm.IsHidden(rawSsid)

		bssid, err := ap.GetProperty(nm.ApHwAddress)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointBssid property: %w", err)
		}

		strength, err := ap.GetProperty(nm.ApStrength)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointStrength property: %w", err)
		}

		flagsRaw, err := ap.GetProperty(nm.ApFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointFlags property: %w", err)
		}
		wpaFlags, err := ap.GetProperty(nm.ApWpaFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointWpaFlags property: %w", err)
		}
		rsnFlags, err := ap.GetProperty(nm.ApRsnFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointRsnFlags property: %w", err)
		}

		flags := flagsRaw.Value().(uint32)
		wpa := wpaFlags.Value().(uint32)
		rsn := rsnFlags.Value().(uint32)
		securityType := nm.DetermineSecurityType(flags, wpa, rsn)

		networks = append(networks, nm.AccessPoint{
			SSID:         ssid,
			BSSID:        bssid.Value().(string),
			SecurityType: securityType,

			DevicePath: wifiDevicePath,
			APPath:     accessPointPath,

			Strength: strength.Value().(uint8),

			IsSaved: false,
			Secured: securityType != "open",
			Hidden:  hidden,
			HasWps:  flags&0x00000002 != 0,
		})
	}

	return networks, nil
}

func GetSavedNetworks(c *dbus.Client, available []nm.AccessPoint) ([]nm.AccessPoint, error) {
	visibleAPs := make(map[string]nm.AccessPoint)
	for _, ap := range available {
		if ap.SSID != "" {
			visibleAPs[ap.SSID] = ap
		}
	}

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

// build an AccessPoint from the connection's settings. Pass in the connection path
// of the connection. available is the slice of available networks that we already
// have so that we can cross-reference/match the new connection with an existing one,
// this is done so we can get info on the network that's not available via it's
// settings - such as Strength.
func GetApFromSettings(c *dbus.Client, path godbus.ObjectPath, available []nm.AccessPoint) (*nm.AccessPoint, error) {
	visibleAPs := make(map[string]nm.AccessPoint)
	for _, ap := range available {
		if ap.SSID != "" {
			visibleAPs[ap.SSID] = ap
		}
	}

	// as per NM docs this is the response type from GetSettings.
	var settings map[string]map[string]godbus.Variant

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

	var ssid string = "unknown"
	var bssid string
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
				bssid = bssidSlice[0]
			}
		}
	}

	var secType string = "unknown"
	_, secured := settings["802-11-wireless-security"]
	if secured {
		if keyMgmtVal, ok := settings["802-11-wireless-security"]["key-mgmt"]; ok {
			secType = keyMgmtVal.Value().(string)
		}
	}

	var strength uint8 = 0
	var activeAPPath godbus.ObjectPath
	var devicePath godbus.ObjectPath

	if activeScan, isNearby := visibleAPs[ssid]; isNearby {
		strength = activeScan.Strength
		bssid = activeScan.BSSID
		activeAPPath = activeScan.APPath
		devicePath = activeScan.DevicePath
	}

	ap := nm.AccessPoint{
		SSID:         ssid,
		BSSID:        bssid,
		SecurityType: secType,

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

func ConnectSaved(client *dbus.Client, network nm.AccessPoint) (godbus.ObjectPath, error) {
	obj, err := nm.ActivateConnection(client, network)

	if err != nil {
		return "", fmt.Errorf("ConnectSaved: %w", err)
	}

	return obj, nil
}

func ConnectSecured(
	client *dbus.Client,
	network *nm.AccessPoint,
	password string,
) error {
	err := nm.AddAndActivateConnection(
		client,
		*network,
		password,
	)

	if err != nil {
		return fmt.Errorf("ConnectSecured: %w", err)
	}

	return nil
}

func ConnectOpen(
	client *dbus.Client,
	network *nm.AccessPoint,
) error {
	err := nm.AddAndActivateConnection(
		client,
		*network,
		"",
	)

	if err != nil {
		return fmt.Errorf("ConnectOpen: %w", err)
	}

	return nil
}

func Disconnect(client *dbus.Client) error {
	ACs, err := nm.GetActiveConnections(client)
	if err != nil {
		return fmt.Errorf("Disconnect: %w", err)
	}
	err = nm.DeactivateConnection(client, ACs)
	if err != nil {
		return fmt.Errorf("Disconnect: %w", err)
	}
	return nil
}
