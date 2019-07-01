package maas

// BlockDevice represents the Block Device endpoint
// TODO: This is only a partial representation of the endpoint.
type BlockDevice struct {
	BlockSize int      `json:"block_size"`
	ID        int      `json:"id"`
	IDPath    string   `json:"id_path"`
	Model     string   `json:"model"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Serial    string   `json:"serial"`
	Size      int      `json:"size"`
	Tags      []string `json:"tags"`
}
