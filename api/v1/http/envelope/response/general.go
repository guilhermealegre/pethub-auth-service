package response

import "github.com/guilhermealegre/go-clean-arch-core-lib/response"

//swagger:model SuccessResponse
type SuccessResponse struct {
	// Success
	Success bool `json:"success"`
}

// swagger:model SwaggerSuccessResponse
type swaggerSuccessResponse struct { //nolint:all
	response.Response
	Data SuccessResponse `json:"data"`
}
