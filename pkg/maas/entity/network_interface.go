package entity

import "net"

// NetworkInterface represents the MaaS Interface endpoint.
type NetworkInterface struct {
	VLAN               VLAN                   `json:"vlan,omitempty"`
	Children           []string               `json:"children,omitempty"`
	Parents            []string               `json:"parents,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	Links              []NetworkInterfaceLink `json:"links,omitempty"`
	Name               string                 `json:"name,omitempty"`
	MACAddress         string                 `json:"mac_address,omitempty"`
	Product            string                 `json:"product,omitempty"`
	FirmwareVersion    string                 `json:"firmware_version,omitempty"`
	SystemID           string                 `json:"system_id,omitempty"`
	Params             interface{}            `json:"params,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Discovered         string                 `json:"discovered,omitempty"`
	Vendor             string                 `json:"vendor,omitempty"`
	ResourceURI        string                 `json:"resource_uri,omitempty"`
	BondXMitHashPolicy string                 `json:"bond_x_mit_hash_policy,omitempty"`
	BondMode           string                 `json:"bond_mode,omitempty"`
	MTU                string                 `json:"mtu,omitempty"`
	EffectiveMTU       int                    `json:"effective_mtu,omitempty"`
	ID                 int                    `json:"id,omitempty"`
	BridgeFD           int                    `json:"bridge_fd,omitempty"`
	BondMIIMon         int                    `json:"bond_mii_mon,omitempty"`
	BondDownDelay      int                    `json:"bond_down_delay,omitempty"`
	BondUpDelay        int                    `json:"bond_up_delay,omitempty"`
	BondLACPRate       int                    `json:"bond_lacp_rate,omitempty"`
	AcceptRA           bool                   `json:"accept_ra,omitempty"`
	Autoconf           bool                   `json:"autoconf,omitempty"`
	Enabled            bool                   `json:"enabled,omitempty"`
	BridgeSTP          bool                   `json:"bridge_stp,omitempty"`
}

// NetworkInterfaceLink is consumed by NetworkInterface{} and should not be used directly.
type NetworkInterfaceLink struct {
	ID        int    `json:"id,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Subnet    Subnet `json:"subnet,omitempty"`
	IPAddress net.IP `json:"ip_address,omitempty"`
}
