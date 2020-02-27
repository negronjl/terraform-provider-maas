package params

// VLAN contains the options for a POST request to the vlans endpoint.
// Only the VID field is required. If Space is empty or the string "undefined",
// the VLAN will be created in the 'undefined' space.
type VLAN struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	VID         int    `json:"vid,omitempty"`
	MTU         int    `json:"mtu,omitempty"`
	Space       string `json:"space,omitempty"`
}
