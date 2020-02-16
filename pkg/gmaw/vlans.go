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

// VLANs provides methods for the vlans operations in the MaaS API.
// This type should be instantiated via NewVLANs(). It fulfills the
// api.VLANs interface.
type VLANs struct {
	c Client
}

// NewVLANs configures a new VLANs.
func NewVLANs(client *gomaasapi.MAASObject) *VLANs {
	c := client.GetSubObject("fabrics")
	return &VLANs{c: Client{&c}}
}

// client returns a Client (ie wrapped MAASOBject) for the VLAN with the given fabric ID
func (v *VLANs) client(fabricID int) Client {
	return v.c.GetSubObject(strconv.Itoa(fabricID)).GetSubObject("vlans")
}

// Get returns information about all of the configured VLANs for <fabric>.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (v *VLANs) Get(fabricID int) (vlans []entity.VLAN, err error) {
	err = v.client(fabricID).Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &vlans)
	})
	return
}

// Post creates a new VLAN and returns information about the new VLAN.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (v *VLANs) Post(fabricID int, p *params.VLAN) (vlan *entity.VLAN, err error) {
	qsp := maas.ToQSP(p)
	vlan = new(entity.VLAN)
	err = v.client(fabricID).Post("", qsp, func(data []byte) error {
		return json.Unmarshal(data, vlan)
	})
	return
}
