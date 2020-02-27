package subnet

import "net"

type IPRange struct {
	Start        net.IP `json:"start,omitempty"`
	End          net.IP `json:"end,omitempty"`
	NumAddresses int    `json:"num_addresses,omitempty"`
}
