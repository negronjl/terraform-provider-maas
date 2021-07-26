package bridge

import (
	"fmt"
	"reflect"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/internal/tfschema"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// NetworkInterface contains methods for connecting maas_interfaces to MaaS Interfaces.
type NetworkInterface struct {
	mo *gomaasapi.MAASObject
}

// NewNetworkInterface creates a new NetworkInterface.
// The parameter should be the metadata passed to the Terraform CRUD functions,
// which should be a *gomaasapi.MAASObject. This function will cast the interface
// received by the Terraform functions to the correct type and store it in the
// NetworkInterface. If the type does not convert, a nil pointer will be stored.
func NewNetworkInterface(m interface{}) *NetworkInterface {
	mo := m.(*gomaasapi.MAASObject)
	return &NetworkInterface{
		mo: mo,
	}
}

// Create a new NetworkInterface in MaaS.
// The sch parameter should be a tfschema type that can be used to create an
// Interface in MaaS. The only compatible type at this time is NetworkInterfacePhysical.
// This method will set the InterfaceID of the type, and expects any attributes
// required to create the Interface to be preset.
// This function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) Create(sch interface{}) error {
	var res *entity.NetworkInterface
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.NetworkInterfacePhysical:
		params := tmpl.Params()
		res, err = gmaw.NewNetworkInterfaces(i.mo).CreatePhysical(tmpl.SystemID, params)
		if err != nil {
			tmpl.InterfaceID = res.ID
		}
	}
	return err
}

// ReadTo updates a tfschema representation of an NetworkInterface to the current state in MaaS.
// The sch parameter should be a tfschema type that represents an Interface in MaaS. This
// function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) ReadTo(sch interface{}) error {
	var res *entity.NetworkInterface
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.NetworkInterfacePhysical:
		if res, err = gmaw.NewNetworkInterface(i.mo).Get(tmpl.SystemID, tmpl.InterfaceID); err != nil {
			return err
		}
		tmpl.Name = res.Name
		tmpl.MACAddress = res.MACAddress
		tmpl.Tags = res.Tags
		tmpl.VLAN = res.VLAN.Name
		tmpl.AcceptRA = res.AcceptRA
		tmpl.Autoconf = res.Autoconf
	}
	return err
}

// UpdateFrom updates the MaaS resource represented by sch.
// The sch parameter should be a tfschema type that represents an Interface in MaaS. This
// function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) UpdateFrom(sch interface{}) error {
	var res *entity.NetworkInterface
	ifc := gmaw.NewNetworkInterface(i.mo)
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.NetworkInterfacePhysical:
		res, err = ifc.Get(tmpl.SystemID, tmpl.InterfaceID)
		if err == nil && !(tmpl.Name == res.Name && tmpl.MACAddress == res.MACAddress &&
			reflect.DeepEqual(tmpl.Tags, res.Tags) && tmpl.VLAN == res.VLAN.Name &&
			tmpl.AcceptRA == res.AcceptRA && tmpl.Autoconf == res.Autoconf) {
			_, err = ifc.Put(tmpl.SystemID, tmpl.InterfaceID, tmpl.Params())
		}
	}
	return err
}

// Delete an Interface in MaaS.
// The sch parameter should be a tfschema type that represents an Interface in MaaS. This
// function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) Delete(sch interface{}) (err error) {
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.NetworkInterfacePhysical:
		err = gmaw.NewNetworkInterface(i.mo).Delete(tmpl.SystemID, tmpl.InterfaceID)
	}
	return
}

// LinkSubnet creates a link between an interface and a subnet.
// This function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) LinkSubnet(sch *tfschema.NetworkInterfaceLink) (err error) {
	_, err = gmaw.NewNetworkInterface(i.mo).LinkSubnet(sch.SystemID, sch.InterfaceID, sch.Params())
	return
}

// UnlinkSubnet removes the link between an interface and a subnet.
// This function will return an error if the MaaS API client returns an error.
func (i *NetworkInterface) UnlinkSubnet(sch *tfschema.NetworkInterfaceLink) (err error) {
	_, err = gmaw.NewNetworkInterface(i.mo).UnlinkSubnet(sch.SystemID, sch.InterfaceID, sch.SubnetID)
	return
}

// ReadLink synchronizes the NetworkInterfaceLink with the current MaaS state.
// This function will return an error if the MaaS API client returns an error, or
// if the link cannot be found in MaaS.
func (i *NetworkInterface) ReadLink(sch *tfschema.NetworkInterfaceLink) error {
	res, err := gmaw.NewNetworkInterface(i.mo).Get(sch.SystemID, sch.InterfaceID)
	if err != nil {
		return err
	}
	for idx := range res.Links {
		if res.Links[idx].Subnet.ID != sch.SubnetID {
			continue
		}
		sch.Mode = res.Links[idx].Mode
		sch.IPAddress = res.Links[idx].IPAddress
	}
	return fmt.Errorf("could not locate link between interface %s.%d and subnet %d",
		sch.SystemID, sch.InterfaceID, sch.SubnetID)
}
