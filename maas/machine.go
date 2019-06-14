package maas

import "encoding/json"

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

// MachineFetcher is the interface that API clients must implement.
type MachineFetcher interface {
	Get(string) ([]byte, error)
	Commission(string, MachineCommissionParams) ([]byte, error)
	Deploy(string, MachineDeployParams) ([]byte, error)
}

type MachineCommissionParams struct{}

type MachineDeployParams struct{}
