package tfschema

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
)

// NetworkInterfaceLink represents a maas_interface_link
type NetworkInterfaceLink struct {
	SystemID       string
	InterfaceID    int
	SubnetID       int
	Mode           string
	IPAddress      net.IP
	Force          bool
	DefaultGateway net.IP
}

// NewNetworkInterfaceLink creates an NetworkInterfaceLink from the Terraform state.
func NewNetworkInterfaceLink(d *schema.ResourceData) *NetworkInterfaceLink {
	var i NetworkInterfaceLink
	if i.SystemID = d.Get("system_id").(string); i.SystemID == "" {
		id := strings.Split(d.Id(), ":")
		i.SystemID = id[0]
		i.InterfaceID, _ = strconv.Atoi(id[1])
		i.SubnetID, _ = strconv.Atoi(id[2])
	} else {
		i.InterfaceID = d.Get("interface_id").(int)
		i.SubnetID = d.Get("subnet_id").(int)
	}
	i.Mode = d.Get("name").(string)
	i.IPAddress = net.ParseIP(d.Get("ip_address").(string))
	i.Force = d.Get("force").(bool)
	i.DefaultGateway = net.ParseIP(d.Get("default_gateway").(string))
	return &i
}

// Params returns a type that can be used to create and delete a MaaS Interface.
func (i *NetworkInterfaceLink) Params() *params.NetworkInterfaceLink {
	return &params.NetworkInterfaceLink{
		Subnet:         i.SubnetID,
		Mode:           i.Mode,
		IPAddress:      i.IPAddress,
		Force:          i.Force,
		DefaultGateway: i.DefaultGateway,
	}
}

// GetID returns "<SystemID>:<InterfaceID>:<SubnetID>" to be used as the Terraform resource ID.
func (i *NetworkInterfaceLink) GetID() (string, error) {
	if i.SystemID == "" {
		return "", fmt.Errorf("SystemID is empty")
	}
	if i.InterfaceID == 0 {
		return "", fmt.Errorf("InterfaceID is zero")
	}
	if i.SubnetID == 0 {
		return "", fmt.Errorf("SubnetID is zero")
	}
	return fmt.Sprintf("%s:%d:%d", i.SystemID, i.InterfaceID, i.SubnetID), nil
}

// UpdateResource updates the Terraform state to reflect the state of the struct.
func (i *NetworkInterfaceLink) UpdateResource(d *schema.ResourceData) (err error) {
	if err = d.Set("mode", i.Mode); err != nil {
		return
	}
	err = d.Set("ip_address", i.IPAddress)
	return
}
