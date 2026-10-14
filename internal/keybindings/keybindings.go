package keybindings

const (
	global = "global"

	// ViewName
	globalView        = global
	wifiSavedView     = "wifi-saved"
	wifiAvailableView = "wifi-available"
	errModalView      = "error-modal"
	passwordView      = "password-modal"

	// Group
	globalGroup        = global
	wifiSavedGroup     = "saved connections"
	wifiAvailableGroup = "available networks"
	navGroup           = "navigation"
	errModalGroup      = "error modal"
	passwordGroup      = "password modal"

	// Keys
	keyEnter = "enter"
	keyEsc   = "esc"
)

type Keybindings struct {
	// NOTE: if you change a bindings ViewName, be sure to update the Focused consts
	// in the root model

	// Global
	FocusedWifiSaved     Binding
	FocusedWifiAvailable Binding
	Cancel               Binding
	OpenKeybindingsModal Binding // NOTE: Cancel is used to close the modal
	Quit                 Binding

	// WiFi
	WifiConnectSaved     Binding
	WifiDisconnect       Binding
	WifiForget           Binding
	WifiAutoConnect      Binding
	WifiScan             Binding
	WifiConnectAvailable Binding

	// Navigation
	Up         Binding
	Down       Binding
	GotoTop    Binding
	GotoBottom Binding

	// Error modal
	ErrorDismiss Binding

	// Password modal
	PasswordVisibility Binding
	PasswordSubmit     Binding
}

type Binding struct {
	// used to match against which View is focused so we can display keybindings for the
	// focused View at the top of the keybindings modal.
	ViewName string
	Group    string // used for grouping in the keybindings modal

	Keys []string

	HelpKey         string // which key to display in the help view
	HelpDescription string // the help text for bubbles key.Binding.Help
	Description     string // the longer full description displayed in the keybindings modal
}

type BindingGroup struct {
	Title    string
	Bindings []Binding
}
