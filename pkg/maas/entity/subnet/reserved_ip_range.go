package subnet

// ReservedIPRange represents an IP range from a Subnet's GetReservedIPRanges()
type ReservedIPRange struct {
	IPRange
	Purpose []string `json:"purpose,omitempty"`
}
