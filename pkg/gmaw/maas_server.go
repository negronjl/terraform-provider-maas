package gmaw

import (
	"fmt"
	"net/url"

	"github.com/juju/gomaasapi"
)

// MAASServer provides methods for the MMAS Server operations in the MaaS API.
// This type should be instantiated via NewMAASServer(). It fulfills the
// api.MAASServer interface.
type MAASServer struct {
	client Client
}

// NewMAASServer configures a new MAASServer.
func NewMAASServer(client *gomaasapi.MAASObject) *MAASServer {
	c := client.GetSubObject("maas")
	return &MAASServer{client: Client{&c}}
}

// Get returns the value of the configuration key <name>.
// This function returns an error if the gomaasapi returns an error.
func (m *MAASServer) Get(name string) (res string, err error) {
	qsp := url.Values{}
	qsp.Set("name", name)
	err = m.client.Get("", qsp, func(data []byte) error {
		res = string(data)
		return nil
	})
	return
}

// Post sets a configuration parameter <name> to <value>.
// This function returns an error if the gomaasapi returns an error.
func (m *MAASServer) Post(name, value string) (err error) {
	qsp := url.Values{}
	qsp.Set("name", name)
	qsp.Set("value", value)
	err = m.client.Post("", qsp, func(data []byte) error {
		if res := string(data); res != "OK" {
			return fmt.Errorf("unexpected server response '%s' (expected 'OK')", res)
		}
		return nil
	})
	return
}
