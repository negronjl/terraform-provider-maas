package gmaw

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// NetworkInterface provides methods for the Interface operations in the MaaS API.
// This type should be instantiated via NewNetworkInterface(). It fulfills the
// api.NetworkInterface interface.
type NetworkInterface struct {
	c Client
}

// NewNetworkInterface configures a new NetworkInterface.
func NewNetworkInterface(client *gomaasapi.MAASObject) *NetworkInterface {
	c := client.GetSubObject("nodes")
	return &NetworkInterface{c: Client{&c}}
}

// client returns a Client with the MAASObject that correlates to the correct endpoint.
func (i *NetworkInterface) client(systemID string, id int) Client {
	return i.c.GetSubObject(systemID).
		GetSubObject("interfaces").
		GetSubObject(strconv.Itoa(id))
}

// Delete the selected interface.
// This function returns an error if the gomaasapi returns an error.
func (i *NetworkInterface) Delete(systemID string, id int) error {
	return i.client(systemID, id).Delete()
}

// Get information about the interface with <id> on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) Get(systemID string, id int) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID, id).Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// AddTag adds an additional tag to the interface.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) AddTag(systemID string, id int, tag string) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	qsp := url.Values{}
	qsp.Add("tag", tag)
	err = i.client(systemID, id).Post("add_tag", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// Disconnect the interface with <id> on <systemID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) Disconnect(systemID string, id int) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID, id).Post("disconnect", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// LinkSubnet links the interface with <id> on <systemID> as described in <params>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) LinkSubnet(systemID string, id int,
	p *params.NetworkInterfaceLink) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	qsp := maas.ToQSP(p)
	err = i.client(systemID, id).Post("link_subnet", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// RemoveTag removes the <tag> tag from the interface
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) RemoveTag(systemID string, id int, tag string) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	qsp := url.Values{}
	qsp.Add("tag", tag)
	err = i.client(systemID, id).Post("remove_tag", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// SetDefaultGateway sets interface <id> to be the default gateway for <systemID>.
// If this interface has more than one subnet with a gateway IP in the same
// IP address family then specifying the ID of the link on this interface is
// required. Set the linkID to 0 to omit this parameter.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) SetDefaultGateway(systemID string, id,
	linkID int) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	qsp := url.Values{}
	if linkID > 0 {
		qsp.Add("link_id", strconv.Itoa(linkID))
	}
	err = i.client(systemID, id).Post("set_default_gateway", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// Unlink subnet removes the link between interface <id> and link <linkID>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) UnlinkSubnet(systemID string, id, linkID int) (ifc *entity.NetworkInterface, err error) {
	ifc = new(entity.NetworkInterface)
	qsp := url.Values{}
	qsp.Add("id", strconv.Itoa(linkID))
	err = i.client(systemID, id).Post("unlink_subnet", qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}

// Put updates the interface configuration with <params>.
// The params argument is one of params.NetworkInterface{Bond,Bridge,Physical,VLAN},
// depending on the type of interface being updated.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (i *NetworkInterface) Put(systemID string, id int, p interface{}) (ifc *entity.NetworkInterface, err error) {
	qsp := maas.ToQSP(p)
	ifc = new(entity.NetworkInterface)
	err = i.client(systemID, id).Put(qsp, func(data []byte) error {
		return json.Unmarshal(data, ifc)
	})
	return
}
