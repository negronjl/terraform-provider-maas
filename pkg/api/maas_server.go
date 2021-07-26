package api

// MAASServer represents the MaaS Server endpoint for changing global configuration settings
type MAASServer interface {
	Get(name string) (value string, err error)
	Post(name, value string) error
}
