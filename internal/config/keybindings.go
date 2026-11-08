package config

import (
	"charm.land/bubbles/v2/key"
)

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
	FocusedInfoPane      Binding
	Cancel               Binding
	OpenHelpMenu         Binding // NOTE: Cancel is used to close the modal
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
	// focused View at the top of the help menu.
	ViewName string
	Group    string // used for grouping in the help menu

	Keys []string

	HelpKey         string // which key to display in the help view
	HelpDescription string // the help text for bubbles key.Binding.Help
	Description     string // the longer full description displayed in the help menu
}

type BindingGroup struct {
	Title    string
	Bindings []Binding
}

var DefaultKeybinds = Keybindings{
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
	FocusedInfoPane: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"3"},
		HelpKey:         "3",
		HelpDescription: "",
		Description:     "Switch focus to info pane",
	},
	Cancel: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{keyEsc},
		HelpKey:         keyEsc,
		HelpDescription: "cancel",
		Description:     "Cancel",
	},
	OpenHelpMenu: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"?"},
		HelpKey:         "?",
		HelpDescription: "Help",
		Description:     "Open help menu",
	},
	Quit: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"q", "ctrl+c"},
		HelpKey:         "q",
		HelpDescription: "quit",
		Description:     "Quit the application",
	},

	// WiFi saved
	WifiConnectSaved: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{keyEnter},
		HelpKey:         keyEnter,
		HelpDescription: "connect",
		Description:     "Connect",
	},
	WifiDisconnect: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"d"},
		HelpKey:         "d",
		HelpDescription: "disconnect",
		Description:     "Disconnect from current WiFi network",
	},
	WifiForget: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"f"},
		HelpKey:         "f",
		HelpDescription: "forget",
		Description:     "Remove the selected WiFi connection",
	},
	WifiAutoConnect: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"a"},
		HelpKey:         "a",
		HelpDescription: "auto-connect",
		Description:     "Toggle auto-connect on/off for the selected WiFi connection",
	},

	// WiFi available
	WifiScan: Binding{
		ViewName:        wifiAvailableView,
		Group:           wifiAvailableGroup,
		Keys:            []string{"s"},
		HelpKey:         "s",
		HelpDescription: "Scan",
		Description:     "Scan for nearby WiFi networks",
	},
	WifiConnectAvailable: Binding{
		ViewName:        wifiAvailableView,
		Group:           wifiAvailableGroup,
		Keys:            []string{keyEnter},
		HelpKey:         keyEnter,
		HelpDescription: "connect",
		Description:     "Connect",
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

// converts a single keybinding to a bubbles keybinding
func (b Binding) ToBubbles() key.Binding {
	return key.NewBinding(
		key.WithKeys(b.Keys...),
		key.WithHelp(b.HelpKey, b.HelpDescription),
	)
}

// converts multiple keybindings to a bubbles keybindings
func ToBubblesBatch(keys []Binding) []key.Binding {
	binds := make([]key.Binding, 0, len(keys))

	for _, k := range keys {
		bind := key.NewBinding(
			key.WithKeys(k.Keys...),
			key.WithHelp(k.HelpKey, k.Description),
		)
		binds = append(binds, bind)
	}

	return binds
}

// Returns all keybindings
func (k Keybindings) All() []Binding {
	return []Binding{
		k.FocusedWifiSaved, k.FocusedWifiAvailable, k.FocusedInfoPane, k.Cancel, k.Quit,
		k.WifiConnectSaved, k.WifiDisconnect, k.WifiForget, k.WifiAutoConnect, k.WifiScan,
		k.WifiConnectAvailable, k.Up, k.Down, k.GotoTop, k.GotoBottom, k.OpenHelpMenu,
		k.ErrorDismiss, k.PasswordVisibility, k.PasswordSubmit,
	}
}

// Get all Keybindings and sort them by Group
func (k Keybindings) AllGrouped() []BindingGroup {
	groups := make([]BindingGroup, 0, 6)
	var global BindingGroup
	var wifiSaved BindingGroup
	var wifiAvailable BindingGroup
	var nav BindingGroup
	var err BindingGroup
	var password BindingGroup

	for _, b := range k.All() {
		if b.Group == globalGroup {
			global = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(globalGroup),
			}
		}
		if b.Group == navGroup {
			nav = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(navGroup),
			}
		}
		if b.Group == wifiSavedGroup {
			wifiSaved = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(wifiSavedGroup),
			}
		}
		if b.Group == wifiAvailableGroup {
			wifiAvailable = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(wifiAvailableGroup),
			}
		}
		if b.Group == errModalGroup {
			err = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(errModalGroup),
			}
		}
		if b.Group == passwordGroup {
			password = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup(passwordGroup),
			}
		}
	}

	groups = append(groups, global)
	groups = append(groups, nav)
	groups = append(groups, wifiSaved)
	groups = append(groups, wifiAvailable)
	groups = append(groups, err)
	groups = append(groups, password)
	return groups
}

// Get keybindings by ViewName
func (k Keybindings) ByView(viewName string) []Binding {
	var binds []Binding
	for _, b := range k.All() {
		if b.ViewName == viewName {
			binds = append(binds, b)
		}
	}
	return binds
}

// Get keybindings by Group
func (k Keybindings) ByGroup(group string) []Binding {
	var binds []Binding
	for _, kb := range k.All() {
		if kb.Group == group {
			binds = append(binds, kb)
		}
	}
	return binds
}
