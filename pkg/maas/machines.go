package maas

// Machines represents the Machines endpoint
type Machines []Machine

// MachinesManager provides locking and management capabilities for Machines
type MachinesManager struct {
	client MachinesFetcher
}

// NewMachineManager creates a new MachinesManager
func NewMachinesManager(client MachinesFetcher) *MachinesManager {
	return &MachinesManager{client: client}
}

// Allocate calls the allocate operation
func (m *MachinesManager) Allocate(params *MachinesAllocateParams) (ma *Machine, err error) {
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
	Tags             []string
	NotTags          []string
	NotInZone        []string
	NotInPool        []string
	Subnets          []string
	NotSubnets       []string
	Storage          []string
	Fabrics          []string
	NotFabrics       []string
	FabricClasses    []string
	NotFabricClasses []string
	Name             string
	SystemID         string
	Arch             string
	Zone             string
	Pool             string
	Pod              string
	NotPod           string
	PodType          string
	NotPodType       string
	Interfaces       string
	AgentName        string
	Comment          string
	CPUCount         int
	Mem              int
	BridgeFD         int
	BridgeAll        bool
	BridgeSTP        bool
	DryRun           bool
	Verbose          bool
}

// MachinesFetcher is the interface that API Clients must implement
type MachinesFetcher interface {
	Allocate(params *MachinesAllocateParams) ([]byte, error)
	Release(systemID []string, comment string) error
}
