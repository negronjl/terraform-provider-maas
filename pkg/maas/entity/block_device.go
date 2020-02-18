package entity

// BlockDevice represents the MaaS BlockDevice endpoint.
type BlockDevice struct {
	BlockSize          int      `json:"block_size,omitempty"`
	ID                 int      `json:"id,omitempty"`
	IDPath             string   `json:"id_path,omitempty"`
	Model              string   `json:"model,omitempty"`
	Name               string   `json:"name,omitempty"`
	Path               string   `json:"path,omitempty"`
	Serial             string   `json:"serial,omitempty"`
	Size               int      `json:"size,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	FirmwareVersion    string   `json:"firmware_version,omitempty"`
	SystemID           string   `json:"system_id,omitempty"`
	AvailableSize      int      `json:"available_size,omitempty"`
	UsedSize           int      `json:"used_size,omitempty"`
	PartitionTableType string   `json:"partition_table_type,omitempty"`
	Partitions         []string `json:"partitions,omitempty"`
	Filesystem         string   `json:"filesystem,omitempty"`
	StoragePool        string   `json:"storage_pool,omitempty"`
	UsedFor            string   `json:"used_for,omitempty"`
	Type               string   `json:"type,omitempty"`
	UUID               string   `json:"uuid,omitempty"`
	ResourceURI        string   `json:"resource_uri,omitempty"`
}
