package api

type MAASServer interface {
	Get(name string) (value string, err error)
	Post(name, value string) error
}
