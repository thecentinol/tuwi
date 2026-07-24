package wifi

import (
	"strconv"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

const (
	secTypeOpen    = "open"
	secTypeUnknown = "unknown"

	ssidHidden = "<hidden>"

	freq24GHz   = "2.4GHz"
	freq5GHz    = "5GHz"
	freq6GHz    = "6GHz"
	freqUnknown = "unknown" //nolint:goconst

	statusConnected   = "Connected"
	statusNearby      = "Nearby"
	statusUnreachable = "Unreachable"
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
		secType = secTypeOpen
	default:
		secType = secTypeUnknown
	}
	return secType
}

func FormatSSID(ssid []byte) string {
	if len(ssid) == 0 {
		return ssidHidden
	}
	return string(ssid)
}

func DetermineFrequency(hertz uint32) string {
	if hertz >= 2401 && hertz <= 2495 {
		return freq24GHz
	}
	if hertz >= 5150 && hertz <= 5895 {
		return freq5GHz
	}
	if hertz >= 5945 && hertz <= 7125 {
		return freq6GHz
	}
	return freqUnknown
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

func FormatStrength(strength uint8) string {
	if strength == 0 {
		return "-"
	}
	return strconv.FormatInt(int64(strength), 10) + "%"
}

// used for saved WiFi networks
func DetermineStatus(isNearby, isActive bool) string {
	switch {
	case isActive:
		return statusConnected

	case isNearby:
		return statusNearby

	default:
		return statusUnreachable
	}
}
