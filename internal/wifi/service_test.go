package wifi

import (
	"github.com/thecentinol/tuwi/internal/dbus"
	"testing"
)

func TestDetermineSecurityType(t *testing.T) {
	tests := []struct {
		name  string
		flags uint32
		wpa   uint32
		rsn   uint32
		want  string
	}{
		{name: "WPA3 Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmtSae, want: "WPA3"},
		{name: "WPA2 Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmtPsk, want: "WPA2"},
		{name: "WPA Network", flags: 0, wpa: dbus.NmSecMgmtPsk, rsn: 0, want: "WPA"},
		{name: "WPA Enterprise Network", flags: 0, wpa: dbus.NmSecMgmt8021, rsn: 0, want: "WPA-ENT"},
		{name: "WPA2 Enterprise Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmt8021, want: "WPA2-ENT"},
		{name: "WPA3 Enterprise Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmtSuiteB192, want: "WPA3-ENT"},
		{name: "OWE Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmtOwe, want: "OWE"},
		{name: "OWE Transition Network", flags: 0, wpa: 0, rsn: dbus.NmSecMgmtOweTm, want: "OWE-TM"},
		{name: "WPA2_WPA3 Network", flags: 0, wpa: dbus.NmSecMgmtPsk, rsn: dbus.NmSecMgmtSae, want: "WPA2/WPA3"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := determineSecurityType(tc.flags, tc.wpa, tc.rsn)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
func FuzzDetermineSecurityType(f *testing.F) {
	f.Add(uint32(0), uint32(0), uint32(0))
	f.Add(uint32(0), uint32(0), uint32(dbus.NmSecMgmtSae))

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
