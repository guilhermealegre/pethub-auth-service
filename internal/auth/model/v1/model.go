package v1

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/guilhermealegre/pethub-auth-service/internal"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

type Model struct {
	app        domain.IApp
	repository domainAuth.IRepository
	streaming  domainAuth.IStreaming
}

func NewModel(
	app domain.IApp,
	repository domainAuth.IRepository,
	streaming domainAuth.IStreaming) domainAuth.IModel {
	return &Model{
		app:        app,
		repository: repository,
		streaming:  streaming,
	}
}

func (m *Model) GetToken(ctx ctxDomain.IContext, loginIdentifier, identifierType, password string) (tokenPair *domainAuth.TokenPair, err error) {

	authDetails, err := m.repository.GetAuthDetails(ctx, m.app.Database().Read(), loginIdentifier, identifierType)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(authDetails.Password, []byte(password))
	if err != nil {
		return nil, err
	}

	tokenJWTDetails := &domainAuth.TokenJWTDetails{
		UserUUID: authDetails.UserUUID,
		Email:    authDetails.Email,
	}

	tokenPair, err = m.generateTokenPair(tokenJWTDetails)
	if err != nil {
		return nil, err
	}

	return tokenPair, err
}

func (m *Model) LoginExternalProviders(ctx ctxDomain.IContext) {
	gothic.BeginAuthHandler(ctx.Response(), ctx.Request())
}

func (m *Model) GetTokenByExternalProviders(ctx ctxDomain.IContext) (*domainAuth.TokenPair, error) {
	gothUser, err := gothic.CompleteUserAuth(ctx.Response(), ctx.Request())
	if err == nil {
		return nil, err
	}

	userDetails := domainAuth.TokenJWTDetails{
		Email: gothUser.Email,
	}

	tokenPair, err := m.generateTokenPair(&userDetails)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (m *Model) Signup(ctx ctxDomain.IContext, loginIdentifier, typeIdentifier string) (err error) {

	if err = m.checkIfLoginIdentifierExists(ctx, loginIdentifier, typeIdentifier); err != nil {
		return err
	}

	confirmationCode := m.generateConfirmationCode()
	if err = m.repository.SaveConfirmationCode(ctx, loginIdentifier, confirmationCode); err != nil {
		return err
	}

	switch typeIdentifier {
	case domainAuth.Email:
		if err = m.streaming.SendEmailSignupConfirmationCode(ctx, loginIdentifier, confirmationCode); err != nil {
			return err
		}
	case domainAuth.PhoneNumber:
		if err = m.streaming.SendSMSSignupConfirmationCode(ctx, loginIdentifier, confirmationCode); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) SignupEmailConfirmation(ctx ctxDomain.IContext, loginIdentifier, code string) (*domainAuth.TokenPair, error) {
	inputCode := strings.TrimSpace(code)
	confirmationCode, err := m.repository.GetConfirmationCode(ctx, loginIdentifier)
	if err != nil {
		return nil, err
	}

	if confirmationCode != inputCode {
		return nil, internal.ErrorTokenConfirmation()
	}

	accessTokenClaims := map[string]interface{}{
		domainAuth.Email: loginIdentifier,
	}

	accessToken, err := m.generateTokenFromClaim(accessTokenClaims, domainAuth.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &domainAuth.TokenPair{
		AccessToken: accessToken,
	}, nil
}

func (m *Model) CreatePassword(ctx ctxDomain.IContext, email, password string) (*domainAuth.TokenPair, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tx, err := m.app.Database().Write().Begin()
	if err != nil {
		return nil, m.app.Logger().DBLog(err)
	}
	defer tx.RollbackUnlessCommitted()

	userUUID, err := m.repository.CreatePassword(ctx, tx, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	_, err = m.streaming.CreateUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	tokenPair, err := m.generateTokenPair(&domainAuth.TokenJWTDetails{
		UserUUID: userUUID,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, m.app.Logger().DBLog(err)
	}

	return tokenPair, nil
}

func (m *Model) Logout(ctx ctxDomain.IContext, idUser int) error {
	return nil
}

func (m *Model) RefreshToken(ctx ctxDomain.IContext, refreshToken string) (*domainAuth.TokenPair, error) {
	return nil, nil
}

func (m *Model) checkIfLoginIdentifierExists(ctx ctxDomain.IContext, loginIdentifier, identifierType string) error {

	switch identifierType {
	case domainAuth.Email:
		userAuthDetails, err := m.repository.GetAuthDetails(ctx, m.app.Database().Read(), loginIdentifier, identifierType)
		if err != nil {
			return err
		}

		if userAuthDetails != nil {
			return internal.ErrorInvalidEmail()
		}

	case domainAuth.PhoneNumber:
		userAuthDetails, err := m.repository.GetAuthDetails(ctx, m.app.Database().Read(), loginIdentifier, identifierType)
		if err != nil {
			return err
		}

		if userAuthDetails != nil {
			return internal.ErrorInvalidPhoneNumber()
		}
	}

	return nil
}

func (m *Model) generateTokenFromClaim(claims map[string]interface{}, tokenTTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims[domainAuth.ExpirationTime] = time.Now().Add(tokenTTL).Unix()
	for key, value := range claims {
		token.Claims.(jwt.MapClaims)[key] = value
	}
	tokenString, err := token.SignedString([]byte(domainAuth.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *Model) generateTokenPair(userDetails *domainAuth.TokenJWTDetails) (*domainAuth.TokenPair, error) {
	accessTokenClaims := map[string]interface{}{
		domainAuth.UserUUID: userDetails.UserUUID,
		domainAuth.Email:    userDetails.Email,
	}
	accessToken, err := m.generateTokenFromClaim(accessTokenClaims, domainAuth.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := map[string]interface{}{
		domainAuth.UserUUID: userDetails.UserUUID,
		domainAuth.Email:    userDetails.Email,
	}
	refreshToken, err := m.generateTokenFromClaim(refreshTokenClaims, domainAuth.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	tokenPair := &domainAuth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func (m *Model) generateConfirmationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
