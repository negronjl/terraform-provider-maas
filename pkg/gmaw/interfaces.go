package gmaw

import (
	"encoding/json"
	"net/url"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// Interfaces provides methods for the Interfaces operations in the MaaS API.
// This type should be instantiated via NewInterfaces(). It fulfills the
// api.Interfaces interface.
type Interfaces struct {
	c Client
}

// NewInterfaces configures a new Interfaces.
func NewInterfaces(client *gomaasapi.MAASObject) *Interfaces {
	c := client.GetSubObject("nodes")
	return &Interfaces{c: Client{&c}}
}

// client returns a Client with the MAASObject that correlates to the correct endpoint.
func (i *Interfaces) client(systemID string) Client {
	return i.c.GetSubObject(systemID).GetSubObject("interfaces")
}

// Get returns information about all of <systemID>'s configured interfaces.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *Interfaces) Get(systemID string) (ifcs []entity.Interface, err error) {
	err = i.client(systemID).Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &ifcs)
	})
	return
}

// CreateBond creates a new bond interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *Interfaces) CreateBond(systemID string, p *params.InterfaceBond) (ifc *entity.Interface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.Interface)
	err = i.client(systemID).Post("create_bond", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreateBridge creates a new bridge interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *Interfaces) CreateBridge(systemID string, p *params.InterfaceBridge) (ifc *entity.Interface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.Interface)
	err = i.client(systemID).Post("create_bridge", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreatePhysical creates a new physical interface on <systemID>'s <params.MACAddress> interface.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *Interfaces) CreatePhysical(systemID string, p *params.InterfacePhysical) (ifc *entity.Interface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.Interface)
	err = i.client(systemID).Post("create_physical", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreateVLAN creates a new VLAN interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *Interfaces) CreateVLAN(systemID string, p *params.InterfaceVLAN) (ifc *entity.Interface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.Interface)
	err = i.client(systemID).Post("create_vlan", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}
