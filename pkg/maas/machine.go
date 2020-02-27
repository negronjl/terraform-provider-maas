package maas

import (
	"encoding/json"
	"sync"

	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// - A proper Caretaker might be necessary here to ease the hackery in Machines.Allocate()
// - There needs to be a way to merge two Machines (ie the one from Allocate) if they have the same SystemID
// - Make an abstract factory so the consumer of this package doesn't have to instantiate each type by hand
// - Convert error responses in gmaw to their "real" analog in juju/errors (why didn't they use this?)

// MACAddress is used by the Machine endpoint
type MACAddress struct {
	MACAddress  string `json:"mac_address"`
	ResourceURI string `json:"resource_uri"`
}

// Machine represents the Machine endpoint
type Machine struct {
	DeployTags             []string             `json:"deploy_tags"`
	Tags                   []string             `json:"tags"`
	IPAddresses            []string             `json:"ip_addresses"`
	MACAddressSet          []MACAddress         `json:"mac_address_set"`
	PhysicalBlockDeviceSet []entity.BlockDevice `json:"physicalblockdevice_set"`
	PXEMac                 []MACAddress         `json:"pxe_mac"`
	Routers                []string             `json:"routers"`
	TagNames               []string             `json:"tag_names"`
	Zone                   []entity.Zone        `json:"zone"`
	Architecture           string               `json:"architecture"`
	BootType               string               `json:"boot_type"`
	DistroSeries           string               `json:"distro_series"`
	Hostname               string               `json:"hostname"`
	DeployHostname         string               `json:"deploy_hostname"`
	OSystem                string               `json:"o_system"`
	Owner                  string               `json:"owner"`
	PowerState             string               `json:"power_state"`
	PowerType              string               `json:"power_type"`
	ResourceURI            string               `json:"resource_uri"`
	SystemID               string               `json:"system_id"`
	UserData               string               `json:"user_data"`
	HWEKernel              string               `json:"hwe_kernel"`
	Comment                string               `json:"comment"`
	CPUCount               int                  `json:"cpu_count"`
	Memory                 int                  `json:"memory"`
	Status                 int                  `json:"status"`
	Storage                int                  `json:"storage"`
	SwapSize               int                  `json:"swap_size"`
	DisableIPv4            bool                 `json:"disable_ipv4"`
	ReleaseErase           bool                 `json:"release_erase"`
	ReleaseEraseSecure     bool                 `json:"release_erase_secure"`
	ReleaseEraseQuick      bool                 `json:"release_erase_quick"`
	Netboot                bool                 `json:"netboot"`
}

// NewMachine converts a MAAS API JSON response into a Golang representation
func NewMachine(data []byte) (m *Machine, err error) {
	m = &Machine{}
	err = json.Unmarshal(data, m)
	return
}

// MachineManager contains functionality for manipulating the Machine endpoint.
// A MachineManager represents a single machine, as does the API endpoint.
type MachineManager struct {
	state  []Machine
	client MachineFetcher
	mutex  sync.RWMutex
}

// NewMachineManager creates a new MachineManager.
// It attempts to fetch the current state of the machine with the given systemID,
// and will return the API client's error if it is not successful. It will also return
// an error from NewMachine if the response cannot successfully be parsed as a Machine.
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
		state:  []Machine{*m},
		client: client,
		mutex:  sync.RWMutex{},
	}, nil
}

// Current returns the most current known state of the machine.
func (m *MachineManager) Current() *Machine {
	return &m.state[len(m.state)-1]
}

func (m *MachineManager) append(current *Machine) {
	m.state = append(m.state, *current)
}

func (m *MachineManager) appendBytes(data []byte) error {
	next, err := NewMachine(data)
	if err == nil {
		m.append(next)
	}
	return err
}

// SystemID returns the machine's systemID.
func (m *MachineManager) SystemID() string {
	return m.Current().SystemID
}

// Commission calls the commission operation on the API.
func (m *MachineManager) Commission(params MachineCommissionParams) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	res, err := m.client.Commission(m.SystemID(), params)
	if err == nil {
		err = m.appendBytes(res)
	}
	return err
}

// Deploy calls the deploy operation on the API.
func (m *MachineManager) Deploy(params *MachineDeployParams) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	res, err := m.client.Deploy(m.SystemID(), params)
	if err == nil {
		err = m.appendBytes(res)
	}
	return err
}

// Lock calls the lock operation on the API.
func (m *MachineManager) Lock(comment string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	res, err := m.client.Lock(m.SystemID(), comment)
	if err == nil {
		err = m.appendBytes(res)
	}
	return err
}

// Update fetches and returns the current state of the machine.
func (m *MachineManager) Update() (ma *Machine, err error) {
	ma, err = m.update()
	if err == nil {
		m.append(ma)
	}
	return
}

func (m *MachineManager) update() (ma *Machine, err error) {
	var res []byte
	res, err = m.client.Get(m.SystemID())
	if err == nil {
		ma, err = NewMachine(res)
	}
	return
}

// MachineFetcher is the interface that API clients must implement.
type MachineFetcher interface {
	Get(string) ([]byte, error)
	Commission(string, MachineCommissionParams) ([]byte, error)
	Deploy(string, *MachineDeployParams) ([]byte, error)
	Lock(string, string) ([]byte, error)
}

// MachineCommissionParams enumerates the parameters for the commission operation
type MachineCommissionParams struct {
	EnableSSH            int
	SkipBMCConfig        int
	SkipNetworking       int
	SkipStorage          int
	CommissioningScripts string
	TestingScripts       string
}

// MachineDeployParams enumerates the parameters for the deploy operation
type MachineDeployParams struct {
	UserData     string
	DistroSeries string
	HWEKernel    string
	AgentName    string
	Comment      string
	BridgeFD     int
	BridgeAll    bool
	BridgeSTP    bool
	InstallRackD bool `json:"install_rackd"`
	InstallKVM   bool `json:"install_kvm"`
}
