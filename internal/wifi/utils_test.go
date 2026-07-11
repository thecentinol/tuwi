package wifi

import (
	"testing"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

func TestDetermineSecurityType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		flags uint32
		wpa   uint32
		rsn   uint32
		want  string
	}{
		{name: "WPA3 Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmtSae, want: "WPA3"},
		{name: "WPA2 Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmtPsk, want: "WPA2"},
		{name: "WPA Network", flags: 0, wpa: nm.NmSecMgmtPsk, rsn: 0, want: "WPA"},
		{name: "WPA Enterprise Network", flags: 0, wpa: nm.NmSecMgmt8021, rsn: 0, want: "WPA-ENT"},
		{name: "WPA2 Enterprise Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmt8021, want: "WPA2-ENT"},
		{name: "WPA3 Enterprise Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmtSuiteB192, want: "WPA3-ENT"},
		{name: "OWE Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmtOwe, want: "OWE"},
		{name: "OWE Transition Network", flags: 0, wpa: 0, rsn: nm.NmSecMgmtOweTm, want: "OWE-TM"},
		{name: "WPA2_WPA3 Network", flags: 0, wpa: nm.NmSecMgmtPsk, rsn: nm.NmSecMgmtSae, want: "WPA2/WPA3"},
		{name: "Open Network with WPS", flags: nm.NmApFlagsWps, wpa: 0, rsn: 0, want: "open"},
		{name: "Open Network no flags", flags: nm.NmApFlagsNone, wpa: 0, rsn: 0, want: "open"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DetermineSecurityType(tc.flags, tc.wpa, tc.rsn)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
func FuzzDetermineSecurityType(f *testing.F) {
	f.Add(uint32(0), uint32(0), uint32(0))
	f.Add(uint32(0), uint32(0), uint32(nm.NmSecMgmtSae))

	f.Fuzz(func(t *testing.T, flags, wpa, rsn uint32) {
		result := DetermineSecurityType(flags, wpa, rsn)
		if result == "" {
			t.Error("should never return empty string")
		}

		if result == "unknown" {
			t.Logf("unknown combination: flags=%d, wpa=%d, rsn=%d", flags, wpa, rsn)
		}

		if flags == 0 && wpa == 0 && rsn == 0 {
			if result != "open" {
				t.Errorf("no flags set but got %q, want open", result)
			}
		}
	})
}

func TestFormatSSID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bytes   []byte
		wantStr string
	}{
		{name: "Hidden Network - nil", bytes: nil, wantStr: "<hidden>"},
		{name: "Hidden Network - empty", bytes: []byte{}, wantStr: "<hidden>"},
		{name: "Visible Network", bytes: []byte("HomeWifi"), wantStr: "HomeWifi"},

		{name: "Visible 1 character ssid", bytes: []byte{65}, wantStr: "A"},
		{name: "Visible with spacing", bytes: []byte("Home Wifi"), wantStr: "Home Wifi"},
		{name: "Visible special characters 1", bytes: []byte("Café_WiFi"), wantStr: "Café_WiFi"},
		{name: "Visible special characters 2", bytes: []byte("@Café_WiFi"), wantStr: "@Café_WiFi"},
		{name: "Max length SSID", bytes: []byte("12345678901234567890123456789012"), wantStr: "12345678901234567890123456789012"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStr := FormatSSID(tc.bytes)
			if gotStr != tc.wantStr {
				t.Errorf(
					"FormatSSID(%v) = got (%v), want (%v)",
					tc.bytes,
					gotStr,
					tc.wantStr,
				)
			}
		})
	}
}

func TestDetermineFrequency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		hertz uint32
		want  string
	}{
		{name: "Unknown Freq with 0 MHz", hertz: uint32(0), want: "unknown"},
		{name: "Uknown Freq with 9999 MHz", hertz: uint32(9999), want: "unknown"},
		{name: "Uknown Freq with 3999 MHz", hertz: uint32(9999), want: "unknown"},
		{name: "Uknown Freq with 7126 MHz", hertz: uint32(9999), want: "unknown"},

		{name: "2.4GHz Freq lower bound", hertz: uint32(2401), want: "2.4GHz"},
		{name: "2.4GHz Freq mid range", hertz: uint32(2442), want: "2.4GHz"},
		{name: "2.4GHz Freq higher bound", hertz: uint32(2495), want: "2.4GHz"},

		{name: "5GHz Freq lower bound", hertz: uint32(5150), want: "5GHz"},
		{name: "5GHz Freq mid range", hertz: uint32(5500), want: "5GHz"},
		{name: "5GHz Freq higher bound", hertz: uint32(5850), want: "5GHz"},

		{name: "6GHz Freq lower bound", hertz: uint32(5945), want: "6GHz"},
		{name: "6GHz Freq mid range", hertz: uint32(6525), want: "6GHz"},
		{name: "6GHz Freq higher bound", hertz: uint32(7125), want: "6GHz"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := DetermineFrequency(tc.hertz)
			if result != tc.want {
				t.Errorf("Incorrect frequency band. Got %v, want %v", result, tc.want)
			}
		})
	}
}

func FuzzDetermineFrequency(f *testing.F) {
	f.Add(uint32(2442))
	f.Add(uint32(5500))
	f.Add(uint32(6525))
	f.Add(uint32(0))
	f.Add(uint32(9999))

	f.Fuzz(func(t *testing.T, hertz uint32) {
		result := DetermineFrequency(hertz)

		if result != "2.4GHz" && result != "5GHz" && result != "6GHz" && result != "unknown" {
			t.Errorf("invalid string: %s", result)
		}
	})
}

func TestDetermineChannel(t *testing.T) {
	tests := []struct {
		name string
		freq uint32
		want uint32
	}{
		// 2.4 GHz
		{"channel 1", 2412, 1},
		{"channel 6", 2437, 6},
		{"channel 11", 2462, 11},
		{"channel 13", 2472, 13},
		{"channel 14", 2484, 14},

		// 5 GHz
		{"channel 32", 5160, 32},
		{"channel 36", 5180, 36},
		{"channel 40", 5200, 40},
		{"channel 44", 5220, 44},
		{"channel 48", 5240, 48},
		{"channel 149", 5745, 149},
		{"channel 161", 5805, 161},
		{"channel 171", 5855, 171},

		// 6 GHz
		{"6ghz channel 2", 5935, 2},
		{"6ghz channel 1", 5955, 1},
		{"6ghz channel 5", 5975, 5},
		{"6ghz channel 233", 7115, 233},

		// Invalid / unknown
		{"below 2.4ghz", 2400, 0},
		{"between 2.4 and 5ghz", 5000, 0},
		{"between 5 and 6ghz", 5900, 0},
		{"above 6ghz", 7200, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DetermineChannel(tc.freq)
			if got != tc.want {
				t.Errorf("DetermineChannel(%d) = %d, want %d", tc.freq, got, tc.want)
			}
		})
	}
}

func TestDetermineStatus(t *testing.T) {
	tests := []struct {
		name       string
		nearbyBool bool
		activeBool bool
		want       string
	}{
		{"Nearby and Active", true, true, "Connected"},
		{"Active only", false, true, "Connected"},
		{"Nearby only", true, false, "Nearby"},
		{"Not nearby or active", false, false, "Unreachable"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DetermineStatus(tc.nearbyBool, tc.activeBool)
			if got != tc.want {
				t.Errorf("DetermineStatus(%v, %v) got: %v, want: %v", tc.nearbyBool, tc.activeBool, got, tc.want)
			}
		})
	}
}
