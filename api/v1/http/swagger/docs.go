/*
	 User Service

	 # User Service API

	 Schemes: http, https
	 BasePath: /v1
	 Version: 1.0

	 Consumes:
	 - application/json

	 Produces:
	 - application/json

	 SecurityDefinitions:
		Bearer:
		  type: apiKey
		  name: Authorization
		  in: header

	 swagger:meta
*/
package swagger

import (
	_ "github.com/guilhermealegre/pethub-auth-service/internal/alive/controller/v1"   // alive controller
	_ "github.com/guilhermealegre/pethub-auth-service/internal/auth/controller/v1"    // auth controller
	_ "github.com/guilhermealegre/pethub-auth-service/internal/swagger/controller/v1" // swagger controller
)
