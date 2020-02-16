package bridge

import (
	"reflect"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/internal/tfschema"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// Interface contains methods for connecting maas_interfaces to MaaS Interfaces.
type Interface struct {
	mo *gomaasapi.MAASObject
}

// NewInterface creates a new Interface.
// The parameter should be the metadata passed to the Terraform CRUD functions,
// which should be a *gomaasapi.MAASObject. This function will cast the interface
// received by the Terraform functions to the correct type and store it in the
// Interface. If the type does not convert, a nil pointer will be stored.
func NewInterface(m interface{}) *Interface {
	mo := m.(*gomaasapi.MAASObject)
	return &Interface{
		mo: mo,
	}
}

// Create a new Interface in MaaS.
// The sch parameter should be a tfschema type that can be used to create an
// Interface in MaaS. The only compatible type at this time is InterfacePhysical.
// This method will set the InterfaceID of the type, and expects any attributes
// required to create the Interface to be preset.
// This function will return an error if the MaaS API client returns an error.
func (i *Interface) Create(sch interface{}) error {
	var res *entity.Interface
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.InterfacePhysical:
		params := tmpl.Params()
		res, err = gmaw.NewInterfaces(i.mo).CreatePhysical(tmpl.SystemID, params)
		if err != nil {
			tmpl.InterfaceID = res.ID
		}
	}
	return err
}

// ReadTo updates a tfschema representation of an Interface to the current state in MaaS.
// The sch parameter should be a tfschema type that represents an Interface in MaaS. This
// function will return an error if the MaaS API client returns an error.
func (i *Interface) ReadTo(sch interface{}) error {
	var res *entity.Interface
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.InterfacePhysical:
		if res, err = gmaw.NewInterface(i.mo).Get(tmpl.SystemID, tmpl.InterfaceID); err != nil {
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
func (i *Interface) UpdateFrom(sch interface{}) error {
	var res *entity.Interface
	ifc := gmaw.NewInterface(i.mo)
	var err error
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.InterfacePhysical:
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
func (i *Interface) Delete(sch interface{}) (err error) {
	switch tmpl := sch.(type) { // nolint: gocritic
	case *tfschema.InterfacePhysical:
		err = gmaw.NewInterface(i.mo).Delete(tmpl.SystemID, tmpl.InterfaceID)
	}
	return
}
