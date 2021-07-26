package subnet

import "net"

// IPRange represents an IP range from a Subnet's GetUnreservedIPRanges()
type IPRange struct {
	Start        net.IP `json:"start,omitempty"`
	End          net.IP `json:"end,omitempty"`
	NumAddresses int    `json:"num_addresses,omitempty"`
}
