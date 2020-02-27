package gmaw

import (
	"encoding/json"
	"net/url"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// Subnets provides methods for the Subnets operations in the MaaS API.
// This type should be instantiated via NewSubnets(). It fulfills the
// api.Subnets interface.
type Subnets struct {
	client Client
}

// NewSubnets configures a new Subnets.
func NewSubnets(client *gomaasapi.MAASObject) *Subnets {
	c := client.GetSubObject("subnets")
	return &Subnets{client: Client{&c}}
}

// Get returns information about all of the configured subnets.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnets) Get() (subnets []entity.Subnet, err error) {
	err = s.client.Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &subnets)
	})
	return
}

// Post creates a new subnet and returns information about the new subnet.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnets) Post(p *params.Subnet) (subnet *entity.Subnet, err error) {
	qsp := maas.ToQSP(p)
	subnet = new(entity.Subnet)
	err = s.client.Post("", qsp, func(data []byte) error {
		return json.Unmarshal(data, subnet)
	})
	return
}
