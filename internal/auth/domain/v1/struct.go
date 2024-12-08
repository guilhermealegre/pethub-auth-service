package v1

import (
	"github.com/google/uuid"
	"time"
)

// swagger:model AdditionalConfigType
type AdditionalConfigType struct {
	UserCacheTtl time.Duration `yaml:"userCacheTtl"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Auth struct {
	Email           string `json:"email"`
	CodePhoneNumber string `json:"code_phone_number"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
}

// swagger:model UserAuthDetails
type UserAuthDetails struct {
	UserUUID        uuid.UUID `json:"user_uuid"`
	Email           string    `json:"email"`
	CodePhoneNumber string    `json:"code_phone_number"`
	PhoneNumber     string    `json:"phone_number"`
	Password        []byte    `json:"password"`
}

type TokenJWTDetails struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Email    string    `json:"email"`
}
