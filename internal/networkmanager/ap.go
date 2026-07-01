package networkmanager

func DetermineSecurityType(flags, wpaFlags, rsnFlags uint32) string {
	var secType string
	switch {
	case rsnFlags&NmSecMgmtSae != 0 && wpaFlags&NmSecMgmtPsk != 0:
		secType = "WPA2/WPA3"
	case rsnFlags&NmSecMgmtSae != 0:
		secType = "WPA3"
	case rsnFlags&NmSecMgmtPsk != 0:
		secType = "WPA2"
	case wpaFlags&NmSecMgmtPsk != 0:
		secType = "WPA"
	case wpaFlags&NmSecMgmt8021 != 0:
		secType = "WPA-ENT"
	case rsnFlags&NmSecMgmt8021 != 0:
		secType = "WPA2-ENT"
	case rsnFlags&NmSecMgmtSuiteB192 != 0:
		secType = "WPA3-ENT"
	case rsnFlags&NmSecMgmtOwe != 0:
		secType = "OWE"
	case rsnFlags&NmSecMgmtOweTm != 0:
		secType = "OWE-TM"
	case flags&NmApFlagsPrivacy != 0:
		secType = "wep"
	case flags&NmApFlagsPrivacy == 0 && wpaFlags == 0 && rsnFlags == 0:
		secType = "open"
	default:
		secType = "unknown"
	}
	return secType
}

func IsHidden(ssidBytes []byte) (bool, string) {
	if len(ssidBytes) == 0 {
		return true, "<hidden>"
	}
	return false, string(ssidBytes)
}
