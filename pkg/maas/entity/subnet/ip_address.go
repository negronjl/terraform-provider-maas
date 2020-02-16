package subnet

import (
	"net"
)

type IPAddress struct {
	IP          net.IP      `json:"ip,omitempty"`
	AllocType   int         `json:"alloc_type,omitempty"`
	Created     string      `json:"created,omitempty"`
	Updated     string      `json:"updated,omitempty"`
	NodeSummary NodeSummary `json:"node_summary,omitempty"`
	User        string      `json:"user,omitempty"`
}

type NodeSummary struct {
	SystemID    string `json:"system_id,omitempty"`
	NodeType    int    `json:"node_type,omitempty"`
	FQDN        string `json:"fqdn,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	IsContainer bool   `json:"is_container,omitempty"`
	Via         string `json:"via,omitempty"`
}
