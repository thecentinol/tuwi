package networkmanager

import (
	"errors"
	"fmt"
	"time"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

const (
	deviceWirelessIface = deviceIface + ".Wireless"
	accessPointIface    = baseIface + ".AccessPoint"

	// Methods
	deviceWirelessGetAllAccessPoints = deviceWirelessIface + ".GetAllAccessPoints"
	deviceWirelessRequestScan        = deviceWirelessIface + ".RequestScan"

	// Signals
	deviceWirelessAccessPointAdded   = deviceWirelessIface + ".AccessPointAdded"
	deviceWirelessAccessPointRemoved = deviceWirelessIface + ".AccessPointRemoved"

	// Properties
	accessPointFlags      = accessPointIface + ".Flags"
	accessPointWpaFlags   = accessPointIface + ".WpaFlags"
	accessPointRsnFlags   = accessPointIface + ".RsnFlags"
	accessPointSsid       = accessPointIface + ".Ssid"
	accessPointFrequency  = accessPointIface + ".Frequency"
	accessPointHwAddress  = accessPointIface + ".HwAddress"
	accessPointMode       = accessPointIface + ".Mode"
	accessPointMaxBitrate = accessPointIface + ".MaxBitrate"
	accessPointBandwidth  = accessPointIface + ".Bandwidth"
	accessPointStrength   = accessPointIface + ".Strength"
	accessPointLastSeen   = accessPointIface + ".LastSeen"

	// NM80211ApFlags
	NmApFlagsNone    = 0x00000000
	NmApFlagsPrivacy = 0x00000001
	NmApFlagsWps     = 0x00000002
	NmApFlagsWpsPbc  = 0x00000004
	NmApFlagsWpsPin  = 0x00000008

	// NM80211ApSecurityFlags
	NmApSecMgmtPsk       = 0x00000100
	NmApSecMgmt8021      = 0x00000200
	NmApSecMgmtSae       = 0x00000400
	NmApSecMgmtOwe       = 0x00000800
	NmApSecMgmtOweTm     = 0x00001000
	NmApSecMgmtSuiteB192 = 0x00002000 // WPA3 Enterprise Suite-B

	// NM80211Mode
	NmModeUnknown = 0
	NmModeAdhoc   = 1
	NmModeInfra   = 2
	NmModeAp      = 3
	NmModeMesh    = 4
)

var DeviceWireless deviceWireless

type deviceWireless struct{}

type AccessPointProperties struct {
	Flags    uint32
	WpaFlags uint32
	RsnFlags uint32

	SSID      []byte
	Frequency uint32
	BSSID     string // HwAddress

	Mode       uint32
	MaxBitrate uint32
	Bandwidth  uint32
	Strength   uint8
	LastSeen   int32

	DevicePath godbus.ObjectPath
	APPath     godbus.ObjectPath
}

// GetAllAccessPoints gets the list of all access points, including hidden ones for
// which the SSID is not yet known.
func (deviceWireless) GetAllAccessPoints(c *dbus.Client, devicePath godbus.ObjectPath,
) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(
		baseIface,
		devicePath,
	)

	call := obj.Call(deviceWirelessGetAllAccessPoints, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetAccessPoints: %w", call.Err)
	}

	var accessPoints []godbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, fmt.Errorf("GetAccessPoints store: %w", err)
	}
	return accessPoints, nil
}

// GetAccessPointProperties gets the properties of an access point
func (deviceWireless) GetAccessPointProperties(
	c *dbus.Client,
	wifiDevicePath godbus.ObjectPath,
	apPath godbus.ObjectPath,
) (*AccessPointProperties, error) {
	obj := c.Conn.Object(baseIface, apPath)

	flags, err := obj.GetProperty(accessPointFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Flags property: %w", err)
	}
	wpaFlags, err := obj.GetProperty(accessPointWpaFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: WpaFlags property: %w", err)
	}
	rsnFlags, err := obj.GetProperty(accessPointRsnFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: RsnFlags property: %w", err)
	}

	ssid, err := obj.GetProperty(accessPointSsid)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: SSID property: %w", err)
	}

	freq, err := obj.GetProperty(accessPointFrequency)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Frequency property: %w", err)
	}

	bssid, err := obj.GetProperty(accessPointHwAddress)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: HwAddress property: %w", err)
	}

	mode, err := obj.GetProperty(accessPointMode)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Mode property: %w", err)
	}

	maxBitate, err := obj.GetProperty(accessPointMaxBitrate)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: MaxBitrate property: %w", err)
	}

	bandwidth, err := obj.GetProperty(accessPointBandwidth)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Bandwidth property: %w", err)
	}

	strength, err := obj.GetProperty(accessPointStrength)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Strength property: %w", err)
	}

	lastSeen, err := obj.GetProperty(accessPointLastSeen)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: LastSeen property: %w", err)
	}

	ap := &AccessPointProperties{
		Flags:    flags.Value().(uint32),
		WpaFlags: wpaFlags.Value().(uint32),
		RsnFlags: rsnFlags.Value().(uint32),

		SSID:      ssid.Value().([]byte),
		Frequency: freq.Value().(uint32),
		BSSID:     bssid.Value().(string),

		Mode:       mode.Value().(uint32),
		MaxBitrate: maxBitate.Value().(uint32),
		Bandwidth:  bandwidth.Value().(uint32),
		Strength:   strength.Value().(uint8),
		LastSeen:   lastSeen.Value().(int32),

		DevicePath: wifiDevicePath,
		APPath:     apPath,
	}

	return ap, nil
}

// RequestScan requests the Wi-Fi device to scan
func (deviceWireless) RequestScan(c *dbus.Client, devicePath godbus.ObjectPath) (err error) {
	device := c.Conn.Object(
		baseIface,
		devicePath,
	)
	call := device.Call(deviceWirelessRequestScan, 0, map[string]godbus.Variant{})
	if call.Err != nil {
		return fmt.Errorf("RequestScan: %w", call.Err)
	}

	// channel for listening to LastScan property
	ch := make(chan *godbus.Signal, 1)
	c.Conn.Signal(ch) // register the channel to receive all signal msgs

	errAdd := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(devicePath),
	)

	if errAdd != nil {
		return fmt.Errorf("RequestScan AddMatchSignal: %w", err)
	}

	defer func() {
		errRm := c.Conn.RemoveMatchSignal(
			godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
			godbus.WithMatchMember("PropertiesChanged"),
			godbus.WithMatchObjectPath(devicePath),
		)
		if errRm != nil {
			err = errors.Join(err, fmt.Errorf("RequestScan RemoveMatchSignal: %w", err))
		}
		c.Conn.RemoveSignal(ch) // deregister the channel
	}()

	select {
	case <-ch:
		// scan complete, continue
	case <-time.After(10 * time.Second):
		return fmt.Errorf("RequestScan: timed out waiting for scan completion")
	}

	return nil
}
