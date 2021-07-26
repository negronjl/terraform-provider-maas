package params

import "net"

// Subnet contains the parameters for the POST operation on the Subnets endpoint.
type Subnet struct {
	CIDR        string   `json:"cidr,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	VLAN        string   `json:"vlan,omitempty"`
	Fabric      string   `json:"fabric,omitempty"`
	VID         int      `json:"vid,omitempty"`
	Space       string   `json:"space,omitempty"`
	GatewayIP   net.IP   `json:"gateway_ip,omitempty"`
	RDNSMode    int      `json:"rdns_mode,omitempty"`
	AllowDNS    bool     `json:"allow_dns,omitempty"`
	AllowProxy  bool     `json:"allow_proxy,omitempty"`
	DNSServers  []string `json:"dns_servers,omitempty"`
	Managed     int      `json:"managed,omitempty"`
}
