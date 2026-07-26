package keybindings

const (
	// ViewName
	globalView        = "global"
	wifiSavedView     = "wifi-saved"
	wifiAvailableView = "wifi-available"
	errModalView      = "error-modal"
	passwordView      = "password-modal"

	// Group
	globalGroup   = "global"
	wifiGroup     = "wifi"
	navGroup      = "navigation"
	errModalGroup = "error modal"
	passwordGroup = "password modal"

	// Keys
	keyEnter = "enter"
	keyEsc   = "esc"
)

func DefaultKeybindings() Keybindings {
	return Keybindings{
		// Global
		FocusedWifiSaved: Binding{
			ViewName:        globalView,
			Group:           globalGroup,
			Keys:            []string{"1"},
			HelpKey:         "1",
			HelpDescription: "",
			Description:     "Switch focus to saved wifi connections pane",
		},
		FocusedWifiAvailable: Binding{
			ViewName:        globalView,
			Group:           globalGroup,
			Keys:            []string{"2"},
			HelpKey:         "2",
			HelpDescription: "",
			Description:     "Switch focus to available wifi networks pane",
		},
		Cancel: Binding{
			ViewName:        globalView,
			Group:           globalGroup,
			Keys:            []string{keyEsc},
			HelpKey:         keyEsc,
			HelpDescription: "cancel",
			Description:     "Cancel",
		},
		OpenKeybindingsModal: Binding{
			ViewName:        globalView,
			Group:           globalGroup,
			Keys:            []string{"?"},
			HelpKey:         "?",
			HelpDescription: "keybindings",
			Description:     "Open keybindings menu",
		},
		Quit: Binding{
			ViewName:        globalView,
			Group:           globalGroup,
			Keys:            []string{"q", "ctrl+c"},
			HelpKey:         "q",
			HelpDescription: "quit",
			Description:     "Quit the application",
		},

		// WiFi
		WifiConnectSaved: Binding{
			ViewName:        wifiSavedView,
			Group:           wifiGroup,
			Keys:            []string{keyEnter},
			HelpKey:         keyEnter,
			HelpDescription: "connect",
			Description:     "Connect to a saved wifi connection",
		},
		WifiDisconnect: Binding{
			ViewName:        wifiSavedView,
			Group:           wifiGroup,
			Keys:            []string{"d"},
			HelpKey:         "d",
			HelpDescription: "disconnect",
			Description:     "Disconnect from current wifi network",
		},
		WifiForget: Binding{
			ViewName:        wifiSavedView,
			Group:           wifiGroup,
			Keys:            []string{"f"},
			HelpKey:         "f",
			HelpDescription: "forget",
			Description:     "Remove the selected wifi connection",
		},
		WifiAutoConnect: Binding{
			ViewName:        wifiSavedView,
			Group:           wifiGroup,
			Keys:            []string{"a"},
			HelpKey:         "a",
			HelpDescription: "auto-connect",
			Description:     "Toggle auto-connect state for the selected wifi connection",
		},
		WifiScan: Binding{
			ViewName:        wifiAvailableView,
			Group:           wifiGroup,
			Keys:            []string{"s"},
			HelpKey:         "s",
			HelpDescription: "Scan",
			Description:     "Scan for nearby wifi networks",
		},
		WifiConnectAvailable: Binding{
			ViewName:        wifiAvailableView,
			Group:           wifiGroup,
			Keys:            []string{keyEnter},
			HelpKey:         keyEnter,
			HelpDescription: "connect",
			Description:     "Connect to a new wifi network",
		},

		// Navigation
		Up: Binding{
			Group:           navGroup,
			Keys:            []string{"k", "up"},
			HelpKey:         "↑",
			HelpDescription: "up",
			Description:     "Move up by one line",
		},
		Down: Binding{
			Group:           navGroup,
			Keys:            []string{"j", "down"},
			HelpKey:         "",
			HelpDescription: "down",
			Description:     "Move down by one line",
		},
		GotoTop: Binding{
			Group:           navGroup,
			Keys:            []string{"t"},
			HelpKey:         "t",
			HelpDescription: "top",
			Description:     "Move/Scroll to top",
		},
		GotoBottom: Binding{
			Group:           navGroup,
			Keys:            []string{"b"},
			HelpKey:         "b",
			HelpDescription: "bottom",
			Description:     "Move/Scroll to bottom",
		},

		// Error modal
		ErrorDismiss: Binding{
			ViewName:        errModalView,
			Group:           errModalGroup,
			Keys:            []string{keyEsc, keyEnter},
			HelpKey:         keyEsc,
			HelpDescription: "dismiss",
			Description:     "Dismiss the error modal",
		},

		// Password modal
		PasswordVisibility: Binding{
			ViewName:        passwordView,
			Group:           passwordGroup,
			Keys:            []string{"ctrl+t"},
			HelpKey:         "ctrl+t",
			HelpDescription: "hide/show password",
			Description:     "Toggle the password input's hidden state",
		},
		PasswordSubmit: Binding{
			ViewName:        passwordView,
			Group:           passwordGroup,
			Keys:            []string{keyEnter},
			HelpKey:         keyEnter,
			HelpDescription: "submit",
			Description:     "Submit password",
		},
	}
}

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
