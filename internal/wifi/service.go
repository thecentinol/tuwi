package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"log"
)

func determineSecurityType(flags, wpaFlags, rsnFlags uint32) string {
	var secType string
	switch {
	case rsnFlags&nm.NmSecMgmtSae != 0 && wpaFlags&nm.NmSecMgmtPsk != 0:
		secType = "WPA2/WPA3"
	case rsnFlags&nm.NmSecMgmtSae != 0:
		secType = "WPA3"
	case rsnFlags&nm.NmSecMgmtPsk != 0:
		secType = "WPA2"
	case wpaFlags&nm.NmSecMgmtPsk != 0:
		secType = "WPA"
	case wpaFlags&nm.NmSecMgmt8021 != 0:
		secType = "WPA-ENT"
	case rsnFlags&nm.NmSecMgmt8021 != 0:
		secType = "WPA2-ENT"
	case rsnFlags&nm.NmSecMgmtSuiteB192 != 0:
		secType = "WPA3-ENT"
	case rsnFlags&nm.NmSecMgmtOwe != 0:
		secType = "OWE"
	case rsnFlags&nm.NmSecMgmtOweTm != 0:
		secType = "OWE-TM"
	case flags&nm.NmApFlagsPrivacy != 0:
		secType = "wep"
	case flags == nm.NmApFlagsNone:
		secType = "open"
	default:
		secType = "unknown"
	}
	return secType
}

func isHidden(ssidBytes []byte) (bool, string) {
	if len(ssidBytes) == 0 {
		return true, "<hidden>"
	}
	return false, string(ssidBytes)
}

// this gets the actual APs + details (ssid, strength, etc)
// for available WiFi networks
func GetAvailableNetworks(c *dbus.Client) ([]nm.AccessPoint, error) {
	wifiDevicePath, err := GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetWifiDevice: %w", err)
	}

	if err := RequestScan(c, wifiDevicePath); err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: RequestScan: %w", err)
	}

	accessPoints, err := GetAccessPoints(c, wifiDevicePath)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: error getting access points: %w", err)
	}

	var networks []nm.AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.Conn.Object(nm.BaseServiceName, accessPointPath)

		ssidBytes, err := ap.GetProperty(nm.AccessPointSsid)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointSsid property: %w", err)
		}
		rawSsid := ssidBytes.Value().([]byte)

		hidden, ssid := isHidden(rawSsid)

		bssid, err := ap.GetProperty(nm.AccessPointBssid)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointBssid property: %w", err)
		}

		strength, err := ap.GetProperty(nm.AccessPointStrength)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointStrength property: %w", err)
		}

		flagsRaw, err := ap.GetProperty(nm.AccessPointFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointFlags property: %w", err)
		}
		wpaFlags, err := ap.GetProperty(nm.AccessPointWpaFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointWpaFlags property: %w", err)
		}
		rsnFlags, err := ap.GetProperty(nm.AccessPointRsnFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointRsnFlags property: %w", err)
		}

		flags := flagsRaw.Value().(uint32)
		wpa := wpaFlags.Value().(uint32)
		rsn := rsnFlags.Value().(uint32)
		securityType := determineSecurityType(flags, wpa, rsn)

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
	err := settingsObj.Call(nm.ListConnections, 0).Store(&connPaths)
	if err != nil {
		return nil, fmt.Errorf("GetSavedNetworks: list connections: %w", err)
	}

	var savedNetworks []nm.AccessPoint

	for _, path := range connPaths {
		conn := c.Conn.Object(nm.BaseServiceName, path)

		// this is the response shape for calling GetSettings,
		// as seen in ...Settings.Connection in the docs
		var settings map[string]map[string]godbus.Variant
		err := conn.Call(nm.GetSettings, 0).Store(&settings)
		if err != nil {
			log.Printf("GetSavedNetworks: failed to get settings for connection: %s: %v", path, err)
			continue
		}

		// filter out non-WiFi connections
		connBlock, ok := settings["connection"]
		if !ok || connBlock["type"].Value().(string) != "802-11-wireless" {
			continue
		}

		// now we can get the SSID & other details
		// NOTE: strength is not available here
		var uuid string
		if v, ok := connBlock["uuid"]; ok {
			uuid = v.Value().(string)
		}

		var ssid string
		if wirelessBlock, ok := settings["802-11-wireless"]; ok {
			if ssidVal, ok := wirelessBlock["ssid"]; ok {
				if ssidBytes, ok := ssidVal.Value().([]byte); ok {
					ssid = string(ssidBytes)
				}
			}
		}

		if ssid == "" {
			continue
		}

		_, secured := settings["802-11-wireless-security"]
		var secType string = "none"
		if secured {
			if keyMgmtVal, ok := settings["802-11-wireless-security"]["key-mgmt"]; ok {
				secType = keyMgmtVal.Value().(string)
			}
		}

		var strength uint8 = 0
		var bssid string = ""
		var activeAPPath godbus.ObjectPath
		var devicePath godbus.ObjectPath

		if activeScan, isNearby := visibleAPs[ssid]; isNearby {
			strength = activeScan.Strength
			bssid = activeScan.BSSID
			activeAPPath = activeScan.APPath
			devicePath = activeScan.DevicePath
		}

		savedNetworks = append(savedNetworks, nm.AccessPoint{
			SSID:           ssid,
			BSSID:          bssid,
			ConnectionUUID: uuid,
			SecurityType:   secType,

			ConnectionPath: path,
			DevicePath:     devicePath,
			APPath:         activeAPPath,

			Strength: strength,

			IsSaved: true,
			Secured: secured,
			Hidden:  false,
		})
	}

	return savedNetworks, nil
}

func GetActiveNetwork(c *dbus.Client) (*nm.AccessPoint, error) {
	var activeNetwork nm.AccessPoint

	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: GetWifiDevice: %w", err)
	}
	device := c.Conn.Object(nm.BaseServiceName, devicePath)
	active, err := device.GetProperty(nm.ActiveAccessPoint)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActiveAccessPoint property: %w", err)
	}
	activePath := active.Value().(godbus.ObjectPath)
	apObject := c.Conn.Object(nm.BaseServiceName, activePath)

	ssid, err := apObject.GetProperty(nm.AccessPointSsid)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActivePointSsid property: %w", err)
	}

	strength, err := apObject.GetProperty(nm.AccessPointStrength)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: AccessPointStrength property: %w", err)
	}

	flags, err := apObject.GetProperty(nm.AccessPointFlags)
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

func ConnectToAvailableSecured(
	client *dbus.Client,
	network *nm.AccessPoint,
	password string,
) error {
	err := AddAndActivateConnection(
		client,
		*network,
		password,
	)

	if err != nil {
		return fmt.Errorf("ConnectToAvailableSecured: AddAndActivateConnection: %w", err)
	}

	return nil
}
