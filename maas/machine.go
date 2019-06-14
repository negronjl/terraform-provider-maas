package maas

import "encoding/json"

// - MachineManager needs a mutex, and
// - A proper Caretaker might be necessary here to ease the hackery in Machines.Allocate()
// - There needs to be a way to merge two Machines (ie the one from Allocate) if they have the same SystemID
// - Need an Update() function to update the machine with the given systemID
// - Might make sense to expose a function to create a machine with the interface mentioned above
// - Make an abstract factory so the consumer of this package doesn't have to instantiate each type by hand

type Machine struct {
	SystemID string `json:"system_id"`
	Hostname string
	// url           string What is this?
	PowerState   string `json:"power_state"`
	CpuCount     int    `json:"cpu_count"`
	Architecture string `optional:"true" forceNew:"true"`
	DistroSeries string `json:"distro_series"`
	Memory       int
	OSystem      string
	Status       int
	TagNames     []string `json:"tag_names"`
	// data          map[string]interface{}
}

func NewMachine(data []byte) (m Machine, err error) {
	err = json.Unmarshal(data, &m)
	return
}

type MachineManager struct {
	state  []Machine
	client MachineFetcher
}

func NewMachineManager(systemID string, client MachineFetcher) (*MachineManager, error) {
	res, err := client.Get(systemID)
	if err != nil {
		return nil, err
	}
	m, err := NewMachine(res)
	if err != nil {
		return nil, err
	}
	return &MachineManager{
		state:  []Machine{m},
		client: client,
	}, nil
}

func (m *MachineManager) Current() Machine {
	return m.state[len(m.state)-1]
}

func (m *MachineManager) append(current Machine) *MachineManager {
	m.state = append(m.state, current)
	return m
}

func (m *MachineManager) appendBytes(data []byte) error {
	next, err := NewMachine(data)
	if err == nil {
		m.append(next)
	}
	return err
}

func (m *MachineManager) SystemID() string {
	return m.Current().SystemID
}

func (m *MachineManager) Commission(params MachineCommissionParams) error {
	res, err := m.client.Commission(m.SystemID(), params)
	if err == nil {
		err = m.appendBytes(res)
	}
	return err
}

func (m *MachineManager) Deploy(params MachineDeployParams) error {
	res, err := m.client.Deploy(m.SystemID(), params)
	if err == nil {
		err = m.appendBytes(res)
	}
	return err
}

// MachineFetcher is the interface that API clients must implement.
type MachineFetcher interface {
	Get(string) ([]byte, error)
	Commission(string, MachineCommissionParams) ([]byte, error)
	Deploy(string, MachineDeployParams) ([]byte, error)
}

type MachineCommissionParams struct{}

type MachineDeployParams struct{}
