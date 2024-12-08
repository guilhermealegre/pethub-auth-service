package v1

import "time"

const (
	Email          = "email"
	PhoneNumber    = "phone_number"
	UserUUID       = "user_uuid"
	ExpirationTime = "exp"

	// token ttl
	AccessTokenTTL  = time.Minute * 15
	RefreshTokenTTL = time.Hour * 24 * 30

	SecretKey = "secret_key_fithub"
)
