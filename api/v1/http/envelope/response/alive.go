package response

import "github.com/guilhermealegre/go-clean-arch-core-lib/response"

// swagger:model SwaggerAliveResponse
type swaggerAliveResponse struct { //nolint:all
	response.Response
	Data AliveResponse `json:"data"`
}

// swagger:model AliveResponse
type AliveResponse struct {
	// Server Name
	ServerName string `json:"serverName"`
	// Port
	Port string `json:"port"`
	// Host Name
	Hostname string `json:"hostname"`
	// Message
	Message string `json:"message"`
}

// swagger:model SwaggerPublicAliveResponse
type swaggerPublicAliveResponse struct { //nolint:all
	response.Response
	Data PublicAliveResponse `json:"data"`
}

// swagger:model PublicAliveResponse
type PublicAliveResponse struct {
	// Name
	Name string `json:"name"`
	// Message
	Message string `json:"message"`
}
