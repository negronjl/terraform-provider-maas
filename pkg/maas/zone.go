package maas

// Zone represents the Zone endpoint
type Zone struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ResourceURI string `json:"resource_uri"`
}
