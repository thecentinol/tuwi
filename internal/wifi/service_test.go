package wifi

import (
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"testing"
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
		{name: "Open Network", flags: 0, wpa: 0, rsn: nm.NmApFlagsNone, want: "open"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := determineSecurityType(tc.flags, tc.wpa, tc.rsn)
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
		result := determineSecurityType(flags, wpa, rsn)
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

func TestIsHidden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		bytes    []byte
		wantBool bool
		wantStr  string
	}{
		{name: "Hidden Network - nil", bytes: nil, wantBool: true, wantStr: "<hidden>"},
		{name: "Hidden Network - empty", bytes: []byte{}, wantBool: true, wantStr: "<hidden>"},
		{name: "Visible Network", bytes: []byte("HomeWifi"), wantBool: false, wantStr: "HomeWifi"},

		{name: "Visible 1 character ssid", bytes: []byte{65}, wantBool: false, wantStr: "A"},
		{name: "Visible with spacing", bytes: []byte("Home Wifi"), wantBool: false, wantStr: "Home Wifi"},
		{name: "Visible special characters 1", bytes: []byte("Café_WiFi"), wantBool: false, wantStr: "Café_WiFi"},
		{name: "Visible special characters 2", bytes: []byte("@Café_WiFi"), wantBool: false, wantStr: "@Café_WiFi"},
		{name: "Max length SSID", bytes: []byte("12345678901234567890123456789012"), wantBool: false, wantStr: "12345678901234567890123456789012"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotBool, gotStr := isHidden(tc.bytes)
			if gotStr != tc.wantStr || gotBool != tc.wantBool {
				t.Errorf(
					"isHidden(%v) = (%v, %q), want (%v, %q)",
					tc.bytes,
					gotBool,
					gotStr,
					tc.wantBool,
					tc.wantStr,
				)
			}
		})
	}
}
