package tfschema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
)

// NetworkInterfacePhysical represents a maas_interface_physical
type NetworkInterfacePhysical struct {
	InterfaceID int
	SystemID    string
	Name        string
	MACAddress  string
	Tags        []string
	VLAN        string
	MTU         int
	AcceptRA    bool
	Autoconf    bool
}

// NewNetworkInterfacePhysical creates an NetworkInterfacePhysical from the Terraform state.
func NewNetworkInterfacePhysical(d *schema.ResourceData) *NetworkInterfacePhysical {
	var i NetworkInterfacePhysical
	if i.SystemID = d.Get("system_id").(string); i.SystemID == "" {
		id := d.Id()
		idx := strings.Index(id, ":")
		i.SystemID = id[:idx]
		i.InterfaceID, _ = strconv.Atoi(id[idx+1:])
	} else {
		i.InterfaceID = d.Get("interface_id").(int)
	}
	i.Name = d.Get("name").(string)
	i.MACAddress = d.Get("mac_address").(string)
	i.VLAN = d.Get("vlan").(string)
	i.MTU = d.Get("mtu").(int)
	i.AcceptRA = d.Get("accept_ra").(bool)
	i.Autoconf = d.Get("autoconf").(bool)
	tags := d.Get("tags").(*schema.Set)
	for _, tag := range tags.List() {
		i.Tags = append(i.Tags, tag.(string))
	}
	return &i
}

// Params returns a type that can be used to create and update a MaaS Interface.
func (i *NetworkInterfacePhysical) Params() *params.NetworkInterfacePhysical {
	return &params.NetworkInterfacePhysical{
		Name:       i.Name,
		MACAddress: i.MACAddress,
		Tags:       i.Tags,
		VLAN:       i.VLAN,
		MTU:        i.MTU,
		AcceptRA:   i.AcceptRA,
		Autoconf:   i.Autoconf,
	}
}

// UpdateResource updates the Terraform state to reflect the state of the struct.
func (i *NetworkInterfacePhysical) UpdateResource(d *schema.ResourceData) (err error) {
	if err = d.Set("name", i.Name); err != nil {
		return
	}
	if err = d.Set("mac_address", i.MACAddress); err != nil {
		return
	}
	if err = d.Set("vlan", i.VLAN); err != nil {
		return
	}
	if err = d.Set("tags", i.Tags); err != nil {
		return
	}
	if err = d.Set("accept_ra", i.AcceptRA); err != nil {
		return
	}
	err = d.Set("autoconf", i.Autoconf)
	return
}

// GetID returns "<SystemID>:<InterfaceID>" to be used as the Terraform resource ID.
func (i *NetworkInterfacePhysical) GetID() (string, error) {
	if i.SystemID == "" {
		return "", fmt.Errorf("SystemID is empty")
	}
	if i.InterfaceID == 0 {
		return "", fmt.Errorf("InterfaceID is zero")
	}
	return fmt.Sprintf("%s:%d", i.SystemID, i.InterfaceID), nil
}
