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
