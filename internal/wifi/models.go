package wifi

type AccessPoint struct {
	Hidden   bool
	SSID     string
	BSSID    string // AKA HwAddress
	Strength uint8
	Secured  bool // derived from privacy flag
	HasWps   bool
}
