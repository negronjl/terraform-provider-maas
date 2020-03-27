package subnet

import "net"

// Statistics represents a Subnet's GetStatistics()
type Statistics struct {
	NumAvailable     int     `json:"num_available,omitempty"`
	LargestAvailable int     `json:"largest_available,omitempty"`
	NumUnavailable   int     `json:"num_unavailable,omitempty"`
	TotalAddresses   int     `json:"total_addresses,omitempty"`
	Usage            float64 `json:"usage,omitempty"`
	UsageString      string  `json:"usage_string,omitempty"`
	AvailableString  string  `json:"available_string,omitempty"`
	FirstAddress     net.IP  `json:"first_address,omitempty"`
	LastAddress      net.IP  `json:"last_address,omitempty"`
	IPVersion        int     `json:"ip_version,omitempty"`
}
