package gmaw

import (
	"encoding/json"
	"net/url"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// RackControllers provides methods for the Rack Controllers operations in the MaaS API.
// This type should be instantiated via NewRackControllers(). It fulfills the
// api.RackControllers interface.
type RackControllers struct {
	client Client
}

// NewRackControllers configures a new RackControllers.
func NewRackControllers(client *gomaasapi.MAASObject) *RackControllers {
	c := client.GetSubObject("rackcontrollers")
	return &RackControllers{client: Client{&c}}
}

// Get returns information about configured rack controllers.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *RackControllers) Get(p *params.RackControllerSearch) (ctrls []entity.RackController, err error) {
	err = s.client.Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &ctrls)
	})
	return
}
