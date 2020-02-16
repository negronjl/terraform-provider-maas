package entity

import "net"

// Subnet represents the MaaS Subnet endpoint.
type Subnet struct {
	Name            string   `json:"name,omitempty"`
	VLAN            VLAN     `json:"vlan,omitempty"`
	CIDR            string   `json:"cidr,omitempty"`
	RDNSMode        int      `json:"rdns_mode,omitempty"`
	GatewayIP       net.IP   `json:"gateway_ip,omitempty"`
	DNSServers      []net.IP `json:"dns_servers,omitempty"`
	AllowDNS        bool     `json:"allow_dns,omitempty"`
	AllowProxy      bool     `json:"allow_proxy,omitempty"`
	ActiveDiscovery bool     `json:"active_discovery,omitempty"`
	Managed         bool     `json:"managed,omitempty"`
	ID              int      `json:"id,omitempty"`
	Space           string   `json:"space,omitempty"`
	ResourceURI     string   `json:"resource_uri,omitempty"`
}
