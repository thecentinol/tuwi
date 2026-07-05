package networkmanager

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func DeleteConnection(c *dbus.Client, conPath string) error {
	obj := c.Conn.Object(BaseServiceName, godbus.ObjectPath(conPath))
	call := obj.Call(CspDelete, 0)

	if call.Err != nil {
		return fmt.Errorf("DeleteConnection: %w", call.Err)
	}

	return nil
}

func GetSettings(c *dbus.Client, path godbus.ObjectPath) (*SavedConnection, error) {
	// as per NM docs this is the response type from GetSettings:
	var settings map[string]map[string]godbus.Variant
	var uuid string
	var autoConnect bool = true
	var ssid string = "unknown"
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

	obj := c.Conn.Object(BaseServiceName, path)
	err := obj.Call(CspGetSettings, 0).Store(&settings)
	if err != nil {
		return nil, fmt.Errorf("GetSettings: store: %w", err)
	}

	// filter out non-WiFi connections
	connBlock, ok := settings["connection"]
	if !ok || connBlock["type"].Value().(string) != "802-11-wireless" {
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

	securityBlock, ok := settings["802-11-wireless-security"]
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

	sc := SavedConnection{
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
		WepKeyType:     wepKeyType,
		WpsMethod:      wpsMethod,
	}

	return &sc, nil
}
