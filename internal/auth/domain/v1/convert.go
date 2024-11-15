package v1

import (
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/request"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/envelope/response"
)

func (t *TokenPair) FromDomainToAPI() *response.AuthResponse {
	if t == nil {
		return nil
	}

	return &response.AuthResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}

func (a *Auth) FromAPIToDomain(reqEmail *request.GetTokenByEmailRequest, reqPhoneNumber *request.GetTokenByPhoneNumberRequest) {
	if (reqEmail == nil && reqPhoneNumber == nil) || a == nil {
		return
	}

	if reqEmail != nil {
		a.Email = reqEmail.Body.Email
		a.Password = reqEmail.Body.Password
	}

	if reqPhoneNumber != nil {
		a.CodePhoneNumber = reqPhoneNumber.Body.CodePhoneNumber
		a.PhoneNumber = reqPhoneNumber.Body.PhoneNumber
		a.Password = reqPhoneNumber.Body.Password

	}
}

func (u *UserDetails) FromAPIToDomain(req *request.EmailSignupConfirmationRequest) *UserDetails {
	if req == nil || u == nil {
		return nil
	}

	return &UserDetails{}

}
