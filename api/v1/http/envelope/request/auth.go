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

		// Password
		// required: true
		Password string `json:"password" validate:"required,gte=6"`
	}
}

type EmailSignupConfirmation struct {
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

type LogoutRequest struct {
	// idUser
	IdUser int `json:"id_user"`
}
