package wifi

import (
	godbus "github.com/godbus/dbus/v5"
	"slices"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

type State struct {
	Saved     []nm.SavedConnection
	Available []nm.AccessPoint

	Nearby []nm.NearbyConnection

	ActiveConnectionPath godbus.ObjectPath
	Scanning             bool
	Error                error
}

func (s *State) SetSaved(connections []nm.SavedConnection) {
	s.Saved = connections
	s.RefreshNearby()
}

func (s *State) AddSaved(conn nm.SavedConnection) {
	s.Saved = append(s.Saved, conn)
	s.RefreshNearby()
}

// when a ConnectionRemovedMsg is received, this function is used to
// remove the selected connection from the SavedConnection (networks) slice.
func (s *State) RemoveSaved(path godbus.ObjectPath) {
	result := slices.DeleteFunc(s.Saved, func(sc nm.SavedConnection) bool {
		return path == sc.ConnectionPath
	})
	s.Saved = result
	s.RefreshNearby()
}

func (s *State) SetAvailable(networks []nm.AccessPoint) {
	s.Available = networks
	s.RefreshNearby()
}

func (s *State) AddAvailable(ap nm.AccessPoint) {
	savedBSSIDs := make(map[string]struct{})
	for _, nearby := range s.Saved {
		for _, bssid := range nearby.BSSIDs {
			savedBSSIDs[bssid] = struct{}{}
		}
	}

	if _, ok := savedBSSIDs[ap.BSSID]; !ok {
		s.Available = append(s.Available, ap)
	}

	s.RefreshNearby()
}

func (s *State) RemoveAvailable(apPath godbus.ObjectPath) {
	result := slices.DeleteFunc(s.Available, func(ap nm.AccessPoint) bool {
		return apPath == ap.APPath
	})
	s.Available = result
	s.RefreshNearby()
}

func (s *State) RefreshNearby() {
	s.Nearby = DisplaySavedConnections(
		s.Saved, s.Available,
	)
}

func (s *State) SetActiveConnection(path godbus.ObjectPath) {
	s.ActiveConnectionPath = path
}

func (s *State) ClearActiveConnection() {
	s.ActiveConnectionPath = ""
}

func (s *State) SetScanning(scanning bool) {
	s.Scanning = scanning
}

func (s *State) SetError(err error) {
	s.Error = err
}
