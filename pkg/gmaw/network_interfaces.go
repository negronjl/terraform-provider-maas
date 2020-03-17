package gmaw

import (
	"encoding/json"
	"net/url"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// NetworkInterfaces provides methods for the Interfaces operations in the MaaS API.
// This type should be instantiated via NewNetworkInterfaces(). It fulfills the
// api.NetworkInterfaces interface.
type NetworkInterfaces struct {
	c Client
}

// NewNetworkInterfaces configures a new NetworkInterfaces.
func NewNetworkInterfaces(client *gomaasapi.MAASObject) *NetworkInterfaces {
	c := client.GetSubObject("nodes")
	return &NetworkInterfaces{c: Client{&c}}
}

// client returns a Client with the MAASObject that correlates to the correct endpoint.
func (i *NetworkInterfaces) client(systemID string) Client {
	return i.c.GetSubObject(systemID).GetSubObject("interfaces")
}

// Get returns information about all of <systemID>'s configured interfaces.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterfaces) Get(systemID string) (ifcs []entity.NetworkInterface, err error) {
	err = i.client(systemID).Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &ifcs)
	})
	return
}

// CreateBond creates a new bond interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterfaces) CreateBond(systemID string,
	p *params.NetworkInterfaceBond) (ifc *entity.NetworkInterface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID).Post("create_bond", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreateBridge creates a new bridge interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterfaces) CreateBridge(systemID string,
	p *params.NetworkInterfaceBridge) (ifc *entity.NetworkInterface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID).Post("create_bridge", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreatePhysical creates a new physical interface on <systemID>'s <params.MACAddress> interface.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterfaces) CreatePhysical(systemID string,
	p *params.NetworkInterfacePhysical) (ifc *entity.NetworkInterface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID).Post("create_physical", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// CreateVLAN creates a new VLAN interface on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterfaces) CreateVLAN(systemID string,
	p *params.NetworkInterfaceVLAN) (ifc *entity.NetworkInterface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID).Post("create_vlan", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}
