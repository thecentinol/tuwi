package wifi

import nm "github.com/thecentinol/tuwi/internal/networkmanager"

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
