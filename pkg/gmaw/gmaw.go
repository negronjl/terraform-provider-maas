/* gmaw is a wrapper (read: adapter) for github.com/juju/gomaasapi to make it
compatible with the client interfaces defined in and expected by the adjacent
maas package. See that package for conventions.

While this package can be used alone, the exported types are designed to be
consumed by their analog in the maas package. This package is just a vehicle
for the maas package to access the MAAS API.

Usage

First, use the provided GetClient() or the gomaasapi equivalent to create a
"global" gomaasapi.MAASObject for consumption by the types defined in this
package. Next, use the New<T> function for the type that correlates to the
endpoint you wish to call to return a usable type, and then start calling
methods on the type to perform operations via your MAAS API.

```
// Get a gomaasapi client
myMAAS, err := gmaw.GetClient("http://example/MAAS/", "supersecr3t", "2.0")
if err != nil {
	log.Fatal(err)
}

// Get the Machine endpoint
myMachineClient := gmaw.NewMachine(myMAAS)

// ...And use it with the maas package
myChassis := maas.NewMachineManager('my_system_id', myMachineClient)
if err := myChassis.Deploy(maas.MachineDeployParams{}); err != nil {
	log.Fatal("Likely got a non-200 response from MAAS API")
}

// Or do it cowboy style
yourMAAS, _ := gmaw.GetClient("http://example/MAAS/", "supersecr3t", "2.0")
res, err := gmaw.NewMachineManager(myMAAS).Deploy('your_sid', maas.MachineDeployParams{})
```
*/
package gmaw

import (
	"net/url"

	"github.com/juju/gomaasapi"
)

// GetClient is a convenience function to create a new gomaasapi client.
// It calls the NewAuthenticatedClient function with the provided credentials,
// and calls NewMAAS() on the returned value to create the resulting MAASObject.
// A non-nil error will be from GetAuthenticatedClient, and the returned client
// can be used with the other types in this package.
func GetClient(apiURL, apiKey, apiVersion string) (*gomaasapi.MAASObject, error) {
	authClient, err := gomaasapi.NewAuthenticatedClient(
		gomaasapi.AddAPIVersionToURL(apiURL, apiVersion), apiKey)
	if err != nil {
		return nil, err
	}
	return gomaasapi.NewMAAS(*authClient), nil
}

type Client struct {
	*gomaasapi.MAASObject
}

func (c Client) Get(op string, params url.Values, f func([]byte) error) error {
	res, err := c.CallGet(op, params)
	if err != nil {
		return err
	}
	data, err := res.GetBytes()
	if err != nil {
		return err
	}
	return f(data)
}

func (c Client) GetSubObject(name string) Client {
	mc := c.MAASObject.GetSubObject(name)
	return Client{&mc}
}

func (c Client) Post(op string, params url.Values, f func([]byte) error) error {
	res, err := c.CallPost(op, params)
	if err != nil {
		return err
	}
	data, err := res.GetBytes()
	if err != nil {
		return err
	}
	return f(data)
}

func (c Client) Put(params url.Values, f func([]byte) error) error {
	res, err := c.Update(params)
	if err != nil {
		return err
	}
	data, err := res.MarshalJSON()
	if err != nil {
		return err
	}
	return f(data)
}
