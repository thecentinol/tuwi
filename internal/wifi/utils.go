package wifi

import (
	"strconv"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

func DetermineSecurityType(flags, wpaFlags, rsnFlags uint32) string {
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
	case flags&nm.NmApFlagsPrivacy == 0 && wpaFlags == 0 && rsnFlags == 0:
		secType = "open"
	default:
		secType = "unknown"
	}
	return secType
}

func FormatSSID(ssid []byte) string {
	if len(ssid) == 0 {
		return "<hidden>"
	}
	return string(ssid)
}

func DetermineFrequency(hertz uint32) string {
	if hertz >= 2401 && hertz <= 2495 {
		return "2.4GHz"
	}
	if hertz >= 5150 && hertz <= 5895 {
		return "5GHz"
	}
	if hertz >= 5945 && hertz <= 7125 {
		return "6GHz"
	}
	return "unknown"
}

func DetermineChannel(freq uint32) uint32 {
	var channel uint32

	if freq >= 2412 && freq <= 2472 {
		channel = (freq - 2407) / 5

	} else if freq == 2484 {
		channel = 14

	} else if freq >= 5160 && freq <= 5855 {
		channel = (freq - 5000) / 5

	} else if freq == 5935 {
		channel = 2

	} else if freq >= 5955 && freq <= 7115 {
		channel = (freq - 5950) / 5
	}
	return channel
}

func FormatChannel(channel uint32) string {
	if channel == 0 {
		return "-"
	}
	return strconv.FormatUint(uint64(channel), 10)
}
