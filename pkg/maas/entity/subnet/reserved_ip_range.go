package subnet

type ReservedIPRange struct {
	IPRange
	Purpose []string `json:"purpose,omitempty"`
}
