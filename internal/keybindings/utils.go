package keybindings

import (
	"charm.land/bubbles/v2/key"
)

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
		k.FocusedWifiSaved, k.FocusedWifiAvailable, k.Cancel, k.Quit, k.WifiConnectSaved,
		k.WifiDisconnect, k.WifiForget, k.WifiAutoConnect, k.WifiScan, k.WifiConnectAvailable,
		k.LineUp, k.LineDown, k.GotoTop, k.GotoBottom, k.OpenKeybindingsModal, k.ErrorDismiss,
		k.PasswordVisibility, k.PasswordSubmit,
	}
}

// Get all Keybindings and sort them by Group
func (k Keybindings) AllGrouped() []BindingGroup {
	groups := make([]BindingGroup, 0, 6)
	var global BindingGroup
	var wifi BindingGroup
	var table BindingGroup
	var kb BindingGroup
	var err BindingGroup
	var password BindingGroup

	for _, b := range k.All() {
		if b.Group == globalGroup {
			global = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("global"),
			}
		}
		if b.Group == wifiGroup {
			wifi = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("wifi"),
			}
		}
		if b.Group == tableGroup {
			table = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("table"),
			}
		}
		if b.Group == keybindsGroup {
			kb = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("keybindings"),
			}
		}
		if b.Group == errModalGroup {
			err = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("error"),
			}
		}
		if b.Group == passwordGroup {
			password = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("password"),
			}
		}
	}

	groups = append(groups, global)
	groups = append(groups, wifi)
	groups = append(groups, table)
	groups = append(groups, kb)
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
