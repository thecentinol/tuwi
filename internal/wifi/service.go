package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
	"github.com/thecentinol/tuwi/internal/models"
)

// this gets the actual APs + details (ssid, strength, etc)
// for available WiFi networks
func GetAvailableNetworks(c *dbus.Client) ([]models.AccessPoint, error) {
	wifiDevicePath, err := dbus.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access points: %v", err)
	}

	if err := dbus.RequestScan(c, wifiDevicePath); err != nil {
		return nil, err
	}

	accessPoints, err := dbus.GetAccessPoints(c, wifiDevicePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access points: %v", err)
	}

	var networks []models.AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.Conn.Object(dbus.BaseServiceName, accessPointPath)

		rawSsid, err := ap.GetProperty(dbus.AccessPointSsid)
		if err != nil {
			return nil, err
		}

		bssid, err := ap.GetProperty(dbus.AccessPointBssid)
		if err != nil {
			return nil, err
		}

		strength, err := ap.GetProperty(dbus.AccessPointStrength)
		if err != nil {
			return nil, err
		}

		flagsRaw, err := ap.GetProperty(dbus.AccessPointFlags)
		if err != nil {
			return nil, err
		}
		wpaFlags, err := ap.GetProperty(dbus.AccessPointWpaFlags)
		if err != nil {
			return nil, err
		}
		rsnFlags, err := ap.GetProperty(dbus.AccessPointRsnFlags)
		if err != nil {
			return nil, err
		}

		flags := flagsRaw.Value().(uint32)
		wpa := wpaFlags.Value().(uint32)
		rsn := rsnFlags.Value().(uint32)

		var secType string
		switch {
		case rsn&dbus.NmSecMgmtSae != 0:
			secType = "sae" // WPA3
		case rsn&dbus.NmSecMgmtPsk != 0:
			secType = "wpa-psk" // WPA2
		case wpa&dbus.NmSecMgmtPsk != 0:
			secType = "wpa-psk" // WPA
		case rsn&dbus.NmSecMgmt8021 != 0 || wpa&dbus.NmSecMgmt8021 != 0:
			secType = "wpa-eap"
		case flags&dbus.NmApFlagsPrivacy != 0:
			secType = "wep"
		default:
			secType = "open"
		}

		ssid := string(rawSsid.Value().([]byte))
		if len(ssid) == 0 {
			ssid = "<hidden>"
		}
		networks = append(networks, models.AccessPoint{
			SSID:         ssid,
			BSSID:        bssid.Value().(string),
			SecurityType: secType,

			DevicePath: wifiDevicePath,
			APPath:     accessPointPath,

			Strength: strength.Value().(uint8),

			IsSaved: false,
			Secured: secType != "open",
			Hidden:  len(ssid) == 0,
			HasWps:  flags&0x00000002 != 0,
		})
	}

	return networks, nil
}

func GetSavedNetworks(c *dbus.Client, available []models.AccessPoint) ([]models.AccessPoint, error) {
	visibleAPs := make(map[string]models.AccessPoint)
	for _, ap := range available {
		if ap.SSID != "" {
			visibleAPs[ap.SSID] = ap
		}
	}

	settingsObj := c.Conn.Object(dbus.BaseServiceName, dbus.SettingsBaseObjPath)

	var connPaths []godbus.ObjectPath
	err := settingsObj.Call(dbus.ListConnections, 0).Store(&connPaths)
	if err != nil {
		return nil, fmt.Errorf("Failed to list connections: %w", err)
	}

	var savedNetworks []models.AccessPoint

	for _, path := range connPaths {
		conn := c.Conn.Object(dbus.BaseServiceName, path)

		// this is the response shape for calling GetSettings,
		// as seen in ...Settings.Connection in the docs
		var settings map[string]map[string]godbus.Variant
		err := conn.Call(dbus.GetSettings, 0).Store(&settings)
		if err != nil {
			// TODO: log error
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

		savedNetworks = append(savedNetworks, models.AccessPoint{
			SSID:           ssid,
			BSSID:          bssid,
			ConnectionUUID: uuid,
			SecurityType:   secType,

			DevicePath: devicePath,
			APPath:     activeAPPath,

			Strength: strength,

			IsSaved: true,
			Secured: secured,
			Hidden:  false,
		})
	}

	return savedNetworks, nil
}

func GetActiveNetwork(c *dbus.Client) (*models.AccessPoint, error) {
	var activeNetwork models.AccessPoint

	devicePath, err := dbus.GetWifiDevice(c)
	if err != nil {
		return nil, err
	}
	device := c.Conn.Object(dbus.BaseServiceName, devicePath)
	active, err := device.GetProperty(dbus.ActiveAccessPoint)
	if err != nil {
		return nil, err
	}
	activePath := active.Value().(godbus.ObjectPath)
	apObject := c.Conn.Object(dbus.BaseServiceName, activePath)

	ssid, err := apObject.GetProperty(dbus.AccessPointSsid)
	if err != nil {
		return nil, err
	}

	strength, err := apObject.GetProperty(dbus.AccessPointStrength)
	if err != nil {
		return nil, err
	}

	flags, err := apObject.GetProperty(dbus.AccessPointFlags)
	if err != nil {
		return nil, err
	}

	ssidBytes := string(ssid.Value().([]byte))
	flagsVal := flags.Value().(uint32)
	activeNetwork = models.AccessPoint{
		SSID:     ssidBytes,
		Strength: strength.Value().(uint8),
		Secured:  flagsVal&0x00000001 != 0,
		HasWps:   flagsVal&0x00000002 != 0,
	}

	return &activeNetwork, nil
}

func ConnectToAvailableSecured(
	client *dbus.Client,
	network *models.AccessPoint,
	password string,
) error {
	err := dbus.AddAndActivateConnection(
		client,
		*network,
		password,
	)

	if err != nil {
		return err
	}

	return nil
}
