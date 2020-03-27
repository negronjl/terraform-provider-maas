/* maas encapsulates a MAAS API client with thread safety, state management,
Go types for resource data, and enhanced functionality for managing API calls.

The types in this package are structured around providing a particular type
of functionality for a given API endpoint: each endpoint will have three
associated types. One type will be named after the endpoint, and maps
directly to the underlying resource: the response to, for example a GET
request to the API endpoint can be unmarshaled into this type. The  New<T>
function for the type will handle this task. Next, each endpoint will have
a <Name>Manager type to handle interacting with the endpoint. The manager
provides the state management, additional safety, and any extra functionality
around API calls. Finally, <Name>Fetcher defines the interface the API client
must expose. More on this below.

Note "endpoint" refers to the headings in the MAAS API Reference
(https://docs.maas.io/2.5/en/api). This means, for example, there is a type
for the Machines (/machines) endpoint as well as for the Machine
(/machines/{system_id}) endpoint.

Usage

In general, a consumer of this package will interact with the Manager types
for each endpoint. Each Manager's New<T> function will require a Fetcher:
the Fetcher will perform the actual API calls, and return the response to
the Manager. The types in the adjacent gmaw package implement the Fetcher
interfaces defined in this package, and the documentation for that package
demonstrates using those types in conjunction with Managers. In general,
it looks something like this:

```
// Get something that fulfills the MachineFetcher interface
machineFetcher := myapiclient.GetMachineFetcher()

// Create the manager
machineManager := maas.NewMachineManager('my_systemid', machineFetcher)

// When we call Deploy() on the manager, the Fetcher makes the actual
// API call and returns the result to the Manager.
machineManager.Deploy(maas.MachineDeployParams{})
```

Defining Endpoint Types

The type that represents the underlying resource is named after the consequent
endpoint, so the "DHCP Snippet" endpoint would correlate to "DHCPSnippet". The
type will be exported and contain exported fields for each field in the
JSON representation of that resource (eg "name", "description", "enabled"),
and the "node" field should be of type Node, where Node performs the same
function for the Node endpoint. As these types serve no purpose other than
to represent the state of a resource at a given time, they have no methods
and no reason to use pointers. Use "json" struct tags where the JSON does
not map directly to a PascalCase Golang name or implement the json.Marshaler
and json.Unmarshaler interfaces. The values of these types' properties should
never be altered; they are only exported to facilitate unmarshaling, and the
types are only exported to save creating a public interface that mimics the type
so other packages can consume the data.

The Manager implements the State pattern, preserving a history of each recorded
change in the underlying resource's state. As consumers of this package will
interact exclusively with the Manager, not the Fetcher, its methods provide
comprehensive access to the API endpoint, but may not map to the endpoint
operations exactly. It provides convenience functions where necessary. The
Manager provides thread safety via Mutexes where necessary, and performs
checks on the current state of the resource before making non-idempotent
API calls. It provides methods to access the current and historical state
of the resource, as well as deltas, and also for determining whether long
running API calls have completed.

The Fetcher is an interface that contains methods that correlate to every
an operation (?op=x) or HTTP verb where there is no operation for a given
endpoint. For example, Machine.Get() correlates to `GET /machines/{system_id}`
and Machine.Deploy() correlates to `POST /machines/{system_id}?op=deploy`. Due
to the extreme variance in parameters between operations, most methods will
accept a type that contains exported fields for each API parameter as a function
parameter. A notable exception is the methods on the Machine type which also
have a systemID string parameter that maps to the variable part of the URL. Each
method will return the raw JSON response from the server as well as an error.
*/
package maas

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// BUG(onwsk8r) JSON does not use unexported fields, and may not convert field names to lower case
// BUG(onwsk8r) ToQSP converts bools to string("true") or string("false"), API wants 0 or 1 <- Verify this

// ToQSP represents a struct as query string parameters.
// Given a struct value as input, it will add each field to a url.Values to be
// printed in the form field=val. Elements of array values are handled via Add().
//
// The behavior of field names is similar to that of the json package: names are
// converted to lower case, the value of a json struct tag will be used if present,
// and the field will be excluded if it has a json tag of "-". Field values are
// printed with fmt.Sprint() to create a string representation; this will work
// properly for simple data types such as int and string.
//
// The function will panic if the input is not a struct, including on a pointer
// to a struct. As this function relies heavily on reflection, it may panic under
// other circumstances as well. It is only expected to work with simple structs
// that represent query string parameters as Go types; other use cases are not
// supported.
func ToQSP(t interface{}) url.Values {
	st := reflect.TypeOf(t)
	sv := reflect.ValueOf(t)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
		sv = sv.Elem()
	}
	qsp := url.Values{}

	for i := 0; i < st.NumField(); i++ {
		// Get the name of the QSP
		key := st.Field(i).Name
		if tag, ok := st.Field(i).Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			key = tag
		}
		key = strings.ToLower(key)

		// Parse out the values
		field := sv.Field(i)
		if field.Kind() == reflect.Array || field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				qsp.Add(key, fmt.Sprint(field.Index(j)))
			}
		} else {
			qsp.Set(key, fmt.Sprint(field))
		}
	}
	return qsp
}
