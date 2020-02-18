package params

// RackControllerSearch narrows down the list in RackController.Get().
// All fields are optional.
type RackControllerSearch struct {
	Hostname   string `json:"hostname,omitempty"`
	MACAddress string `json:"mac_address,omitempty"`
	SystemID   string `json:"id,omitempty"`
	Domain     string `json:"domain,omitempty"`
	Zone       string `json:"zone,omitempty"`
	Pool       string `json:"pool,omitempty"`
	AgentName  string `json:"agent_name,omitempty"`
}
