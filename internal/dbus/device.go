package dbus

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/models"
	"log"
	"time"
)

type NetworkManager struct {
	client *Client
}

func GetDevices(c *Client) ([]dbus.ObjectPath, error) {
	nm := c.conn.Object(baseServiceName, baseObjPath)
	call := nm.Call(baseServiceName+".GetDevices", 0)
	if call.Err != nil {
		return nil, call.Err
	}
	var devices []dbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func GetWifiDevice(c *Client) (dbus.ObjectPath, error) {
	devices, err := GetDevices(c)
	if err != nil {
		log.Fatalf("Error fetching devices: %v", err)
	}

	// finding the wifi device
	for _, v := range devices {
		device := c.conn.Object(baseServiceName, v)
		variant, err := device.GetProperty(deviceType)
		if err != nil {
			return "", err
		}
		if deviceTypeVal, ok := variant.Value().(uint32); ok {
			if deviceTypeVal == 2 {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("No wifi device found")
}

func RequestScan(c *Client) error {
	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return err
	}

	device := c.conn.Object(
		baseServiceName,
		devicePath,
	)
	call := device.Call(deviceWireless+".RequestScan", 0, map[string]dbus.Variant{})
	if call.Err != nil {
		return call.Err
	}

	// channel for listening to LastScan property
	ch := make(chan *dbus.Signal, 1)
	c.conn.Signal(ch) // register the channel to receive all signal msgs

	c.conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
		dbus.WithMatchObjectPath(devicePath),
	)

	select {
	case <-ch:
		// scan complete, continue
	case <-time.After(10 * time.Second):
		return fmt.Errorf("scan timed out")
	}
	return nil
}

// this gets the raw paths for the access points
func GetAccessPoints(c *Client) ([]dbus.ObjectPath, error) {
	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return nil, err
	}

	if err := RequestScan(c); err != nil {
		return nil, err
	}

	nm := c.conn.Object(
		baseServiceName,
		devicePath,
	)

	call := nm.Call(deviceWireless+".GetAllAccessPoints", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var accessPoints []dbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, err
	}
	return accessPoints, nil
}

// this gets the actual APs + details (ssid, strength, etc)
// for available WiFi networks
func GetNetworks(c *Client) ([]models.AccessPoint, error) {
	accessPoints, err := GetAccessPoints(c)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access points: %v", err)
	}

	var networks []models.AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.conn.Object(baseServiceName, accessPointPath)

		ssid, err := ap.GetProperty(accessPointSsid)
		if err != nil {
			return nil, err
		}

		bssid, err := ap.GetProperty(accessPointBssid)
		if err != nil {
			return nil, err
		}

		strength, err := ap.GetProperty(accessPointStrength)
		if err != nil {
			return nil, err
		}

		flags, err := ap.GetProperty(accessPointFlags)
		if err != nil {
			return nil, err
		}

		ssidBytes := string(ssid.Value().([]byte))
		flagsVal := flags.Value().(uint32)
		networks = append(networks, models.AccessPoint{
			Hidden:   ssidBytes == "",
			SSID:     ssidBytes,
			BSSID:    bssid.Value().(string),
			Strength: strength.Value().(uint8),
			Secured:  flagsVal&0x00000001 != 0,
			HasWps:   flagsVal&0x00000002 != 0,
		})
	}

	return networks, nil
}

func GetSavedNetworks(c *Client) ([]models.AccessPoint, error) {
	settings := c.conn.Object(baseServiceName, settingsBaseObjPath)
	call := settings.Call(listConnections, 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var connPaths []dbus.ObjectPath
	if err := call.Store(&connPaths); err != nil {
		return nil, err
	}

	var savedNetworks []models.AccessPoint

	// this is used for cross-referencing to get more details,
	// about the network that we cant get from settings. Such as Strength.
	networks, err := GetNetworks(c)
	if err != nil {
		return nil, err
	}

	for _, path := range connPaths {
		conn := c.conn.Object(baseServiceName, path)
		pathCall := conn.Call(getSettings, 0)

		// this is the response shape for calling GetSettings,
		// as seen in ...Settings.Connection in the docs
		var settings map[string]map[string]dbus.Variant
		if err := pathCall.Store(&settings); err != nil {
			return nil, err
		}

		// filter out non-WiFi connections
		if settings["connection"]["type"].Value().(string) != "802-11-wireless" {
			continue
		}

		// now we can get the SSID & other details
		// NOTE: strength is not available here
		ssid := string(settings["802-11-wireless"]["ssid"].Value().([]byte))
		_, secured := settings["802-11-wireless-security"] // checking if key exists
		// get the full slice of bssids to compare the
		bssids := settings["802-11-wireless"]["seen-bssids"].Value().([]string)

		var strength uint8
		for _, network := range networks {
			for _, bssid := range bssids {
				if network.BSSID == bssid {
					strength = network.Strength
				}
			}
		}

		savedNetworks = append(savedNetworks, models.AccessPoint{
			Hidden:   ssid == "",
			SSID:     ssid,
			BSSID:    bssids[0],
			Strength: strength,
			Secured:  secured,
		})
	}

	return savedNetworks, nil
}

func GetActiveNetwork(c *Client) (*models.AccessPoint, error) {
	var activeNetwork models.AccessPoint

	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return nil, err // NOTE: is return type correct
	}
	device := c.conn.Object(baseServiceName, devicePath)
	active, err := device.GetProperty(activeAccessPoint)
	if err != nil {
		return nil, err
	}
	activePath := active.Value().(dbus.ObjectPath)
	apObject := c.conn.Object(baseServiceName, activePath)

	ssid, err := apObject.GetProperty(accessPointSsid)
	if err != nil {
		return nil, err
	}

	strength, err := apObject.GetProperty(accessPointStrength)
	if err != nil {
		return nil, err
	}

	flags, err := apObject.GetProperty(accessPointFlags)
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
