package maas

// Machines represents the Machines endpoint
type Machines []Machine

// MachinesManager
type MachinesManager struct {
	client MachinesFetcher
}

// NewMachineManager
func NewMachinesManager(client MachinesFetcher) *MachinesManager {
	return &MachinesManager{client: client}
}

// Allocate calls the allocate operation
func (m *MachinesManager) Allocate(params MachinesAllocateParams) (ma Machine, err error) {
	var res []byte
	res, err = m.client.Allocate(params)
	if err == nil {
		ma, err = NewMachine(res)
	}
	return
}

// Release calls the release operation.
func (m *MachinesManager) Release(systemIDs []string, comment string) error {
	return m.client.Release(systemIDs, comment)
}

// MachinesParams enumerates the options for the GET operation
type MachinesParams struct {
	Hostname   []string
	MACAddress []string
	ID         []string
	Domain     string
	Zone       string
	AgentName  string
	Pool       []string
}

// MachinesAllocateParams enumerates the options for the allocate operation.
type MachinesAllocateParams struct {
	Name             string
	SystemID         string
	Arch             string
	CPUCount         int
	Mem              int
	Tags             []string
	NotTags          []string
	Zone             string
	NotInZone        []string
	Pool             string
	NotInPool        []string
	Pod              string
	NotPod           string
	PodType          string
	NotPodType       string
	Subnets          []string
	NotSubnets       []string
	Storage          []string
	Interfaces       string
	Fabrics          []string
	NotFabrics       []string
	FabricClasses    []string
	NotFabricClasses []string
	AgentName        string
	Comment          string
	BridgeAll        bool
	BridgeSTP        bool
	BridgeFD         int
	DryRun           bool
	Verbose          bool
}

// MachinesFetcher is the interface that API Clients must implement
type MachinesFetcher interface {
	Allocate(params MachinesAllocateParams) ([]byte, error)
	Release(systemId []string, comment string) error
}
