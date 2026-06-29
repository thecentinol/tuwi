package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
)

func GetWifiDevice(c *dbus.Client) (godbus.ObjectPath, error) {
	devices, err := GetDevices(c)
	if err != nil {
		return "", fmt.Errorf("GetWifiDevice: %w", err)
	}

	// finding the wifi device
	for _, v := range devices {
		device := c.Conn.Object(BaseServiceName, v)
		variant, err := device.GetProperty(DeviceType)
		if err != nil {
			return "", fmt.Errorf("GetWifiDevice: failed to get DeviceType property %w", err)
		}
		if deviceTypeVal, ok := variant.Value().(uint32); ok {
			if deviceTypeVal == 2 {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("GetWifiDevice: no wifi device found")
}
