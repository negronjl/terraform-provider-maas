package maas

// BlockDevice represents the Block Device endpoint
// NOTE: This is only a partial representation of the endpoint.
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
