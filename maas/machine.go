package maas

// MachineFetcher is the interface that API clients must implement.
type MachineFetcher interface {
	Get(string) ([]byte, error)
	Commission(string, MachineCommissionParams) ([]byte, error)
	Deploy(string, MachineDeployParams) ([]byte, error)
}

type MachineCommissionParams struct{}

type MachineDeployParams struct{}
