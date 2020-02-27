package entity

// Pod represents the MaaS Pod endpoint.
type Pod struct {
	Zone      Zone         `json:"zone,omitempty"`
	Pool      ResourcePool `json:"pool,omitempty"`
	Used      PodResource  `json:"used,omitempty"`
	Available PodResource  `json:"available,omitempty"`
	Total     PodResource  `json:"total,omitempty"`
	Host      struct {
		SystemID   string `json:"system_id,omitempty"`
		Incomplete bool   `json:"__incomplete__,omitempty"`
	} `json:"host,omitempty"`
	StoragePools          []PodStoragePool `json:"storage_pools,omitempty"`
	Architectures         []string         `json:"architectures,omitempty"`
	Tags                  []string         `json:"tags,omitempty"`
	Capabilities          []string         `json:"capabilities,omitempty"`
	Name                  string           `json:"name,omitempty"`
	Type                  string           `json:"type,omitempty"`
	DefaultMACVLANMode    string           `json:"default_macvlan_mode,omitempty"`
	ResourceURI           string           `json:"resource_uri,omitempty"`
	CPUOverCommitRatio    float64          `json:"cpu_over_commit_ratio,omitempty"`
	MemoryOverCommitRatio float64          `json:"memory_over_commit_ratio,omitempty"`
	ID                    int              `json:"id,omitempty"`
}

// PodStoragePool represents the "storage_pools" object in a Pod.
// This type should not be used directly.
type PodStoragePool struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Path      string `json:"path,omitempty"`
	Total     int    `json:"total,omitempty"`
	Used      int    `json:"used,omitempty"`
	Available int    `json:"available,omitempty"`
	Default   bool   `json:"default,omitempty"`
}

// PodResource represents the "used", "available", and "total" objects in a Pod
// This type should not be used directly.
type PodResource struct {
	Cores        int `json:"cores,omitempty"`
	Memory       int `json:"memory,omitempty"`
	LocalStorage int `json:"local_storage,omitempty"`
}
