package provider

// BlockDevice is used by the maas_instance resource
type BlockDevice struct {
	BlockSize int
	ID        int
	IDPath    string
	Model     string
	Name      string
	Path      string
	Serial    string
	Size      int
	Tags      []string
}
