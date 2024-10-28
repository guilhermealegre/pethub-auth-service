package response

import "github.com/guilhermealegre/go-clean-arch-core-lib/response"

// swagger:model SwaggerScannerLoginResponse
type swaggerScannerLoginResponse struct { //nolint:all
	response.Response
	Data AuthResponse `json:"data"`
}

// swagger:model SwaggerCreatePasswordResponse
type swaggerCreatePasswordResponse struct {
	response.Response
	Data AuthResponse `json:"data"`
}

//swagger:model AuthResponse
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
