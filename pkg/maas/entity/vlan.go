package entity

// VLAN represents the MaaS VLAN endpoint.
type VLAN struct {
	VID           int    `json:"vid,omitempty"`
	MTU           int    `json:"mtu,omitempty"`
	DHCPOn        bool   `json:"dhcp_on,omitempty"`
	ExternalDHCP  string `json:"external_dhcp,omitempty"`
	RelayVLAN     int    `json:"relay_vlan,omitempty"`
	FabricID      int    `json:"fabric_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	PrimaryRack   string `json:"primary_rack,omitempty"`
	SecondaryRack string `json:"secondary_rack,omitempty"`
	Space         string `json:"space,omitempty"`
	ID            int    `json:"id,omitempty"`
	Fabric        string `json:"fabric,omitempty"`
	ResourceURI   string `json:"resource_uri,omitempty"`
}
