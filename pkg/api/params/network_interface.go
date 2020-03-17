package params

import "net"

// NetworkInterfaceBond is the parameters for the NetworkInterfaces create_bond POST operation.
type NetworkInterfaceBond struct {
	NetworkInterfacePhysical
	Parents            []int  `json:"parents,omitempty"`
	BondMode           string `json:"bond_mode,omitempty"`
	BondMiimon         int    `json:"bond_miimon,omitempty"`
	BondDownDelay      int    `json:"bond_down_delay,omitempty"`
	BondUpDelay        int    `json:"bond_up_delay,omitempty"`
	BondLACPRate       string `json:"bond_lacp_rate,omitempty"`
	BondXMitHashPolicy string `json:"bond_x_mit_hash_policy,omitempty"`
	BondNumberGratARP  int    `json:"bond_number_grat_arp,omitempty"`
}

// NetworkInterfaceBridge is the parameters for the NetworkInterfaces create_bridge POST operation.
type NetworkInterfaceBridge struct {
	NetworkInterfacePhysical
	Parent    int  `json:"parent,omitempty"`
	BridgeSTP bool `json:"bridge_stp,omitempty"`
	BridgeFD  int  `json:"bridge_fd,omitempty"`
}

// NetworkInterfacePhysical is the parameters for the NetworkInterfaces create_physical POST operation.
type NetworkInterfacePhysical struct {
	Name       string   `json:"name,omitempty"`
	MACAddress string   `json:"mac_address,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	VLAN       string   `json:"vlan,omitempty"`
	MTU        int      `json:"mtu,omitempty"`
	AcceptRA   bool     `json:"accept_ra,omitempty"`
	Autoconf   bool     `json:"autoconf,omitempty"`
}

// NetworkInterfaceVLAN is the parameters for the NetworkInterfaces create_vlan POST operation.
type NetworkInterfaceVLAN struct {
	Tags     []string `json:"tags,omitempty"`
	VLAN     string   `json:"vlan,omitempty"`
	Parent   int      `json:"parent,omitempty"`
	MTU      int      `json:"mtu,omitempty"`
	AcceptRA bool     `json:"accept_ra,omitempty"`
	Autoconf bool     `json:"autoconf,omitempty"`
}

// NetworkInterfaceLink is used with NetworkInterface.LinkSubnet().
// Mode must be one of (AUTO, DHCP, STATIC, LINK_UP). IPAddress is ignored
// unless mode is STATIC, and will be set automatically if empty. Force
// allows LINK_UP to be set when other links exist, allows links between
// different VLANs, and deletes all other links on the interface.
// DefaultGateway is ignored unless Mode is AUTO or STATIC.
// Note: You can parse an IP address into a net.IP via net.ParseIP(string).
type NetworkInterfaceLink struct {
	Mode           string `json:"mode,omitempty"`
	Subnet         int    `json:"subnet,omitempty"`
	IPAddress      net.IP `json:"ip_address,omitempty"`
	Force          bool   `json:"force,omitempty"`
	DefaultGateway net.IP `json:"default_gateway,omitempty"`
}
