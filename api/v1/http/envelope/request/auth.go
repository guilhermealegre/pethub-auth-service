package request

//swagger:parameters ScannerLoginRequest
type GetTokenByEmailRequest struct {

	// Body
	// in: body
	// required: true
	Body struct {
		// Email
		// required: true
		Email string `json:"email" validate:"required,email"`
		// Password
		// required: true
		Password string `json:"password" validate:"required,gte=6"`
	}
}

//swagger:parameters GetTokenByPhoneNumberRequest
type GetTokenByPhoneNumberRequest struct {
	// Body
	// in: body
	// required: true
	Body struct {
		// Phone Number Code
		// required: true
		CodePhoneNumber string `json:"code_phone_number" validate:"required,gte=1"`
		// Phone Number
		// required: true
		PhoneNumber string `json:"phone_number" validate:"required,gte=1"`
		// Password
		// required: true
		Password string `json:"password" validate:"required,gte=6"`
	}
}

//swagger:parameters SignUpEmailRequest
type SignUpEmailRequest struct {
	// Body
	// in: body
	// required: true
	Body struct {
		// Email
		// required: true
		Email string `json:"email" validate:"email"`
	}
}

//swagger:parameters SignUpPhoneNumberRequest
type SignUpPhoneNumberRequest struct {
	// Body
	// in: body
	// required: true
	Body struct {
		// Phone Number
		// required: true
		PhoneNumber string `json:"phone_number" validate:"required, gte=1"`
		// Code Phone Number
		// required: true
		CodePhoneNumber string `json:"code_phone_number" validate:"required, gte=1"`
	}
}

//swagger:parameters EmailSignupConfirmationRequest
type EmailSignupConfirmationRequest struct {
	// Body
	// in: body
	// required: true
	Body struct {
		// Code
		// required: true
		Code string `json:"code" validate:"required,eq=6"`
		// Email
		// required: true
		Email string `json:"email" validate:"required,email"`
	}
}

//swagger:parameters SignUpPhoneNumberConfirmationRequest
type SignUpPhoneNumberConfirmationRequest struct {
	// Body
	// in: body
	// required: true
	Body struct {
		// Code
		// required: true
		Code string `json:"code" validate:"required,eq=6"`
		// Phone Number
		// required: true
		PhoneNumber string `json:"phone_number" validate:"required, gte=1"`
		// Code Phone Number
		// required: true
		CodePhoneNumber string `json:"code_phone_number" validate:"required, gte=1"`
	}
}

//swagger:parameters CreatePasswordRequest
type CreatePasswordRequest struct {
	Body struct {
		// Password
		// required: true
		Password string `json:"password" validate:"required,gte=6"`
		// Confirm Password
		// required: true
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	}
}

//swagger:parameters LogoutRequest
type LogoutRequest struct {
	// idUser
	IdUser int `json:"id_user"`
}
