package v1

// swagger:model Alive
type Alive struct {
	// Server Name
	ServerName string `json:"server_name"`
	// Port
	Port string `json:"port"`
	// Host Name
	Hostname string `json:"hostname"`
	// Message
	Message string `json:"message"`
}

// swagger:model PublicAlive
type PublicAlive struct {
	// Name
	Name string `json:"name"`
	// Message
	Message string `json:"message"`
}
