package v1

import (
	"fmt"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/pethub-auth-service/internal"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
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

func (m *Model) GetToken(ctx domain.IContext, loginIdentifier, identifierType, password string) (tokenPair *domainAuth.TokenPair, err error) {

	userAuthDetails, err := m.repository.GetAuthDetails(ctx, m.app.Database().Read(), loginIdentifier, identifierType)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(userAuthDetails.Password, []byte(password))
	if err != nil {
		return nil, err
	}

	userDetails, err := m.repository.GetUserDetails(ctx, userAuthDetails.IdUser)
	if err != nil {
		return nil, err
	}

	userDetails.Email = userAuthDetails.Email
	tokenPair, err = m.generateTokenPair(userDetails)
	if err != nil {
		return nil, err
	}

	return tokenPair, err
}

func (m *Model) Signup(ctx domain.IContext, loginIdentifier, typeIdentifier string) (err error) {
	if err = m.checkIfLoginIdentifierExists(ctx, loginIdentifier, typeIdentifier); err != nil {
		return err
	}

	confirmationCode := m.generateConfirmationCode()
	err = m.repository.SaveConfirmationCode(ctx, loginIdentifier, confirmationCode)
	if err != nil {
		return err
	}

	switch typeIdentifier {
	case domainAuth.Email:
		go func() { /*
				err := m.modelNotification.SendEmail(
					ctx,
					push_email.EmailSignupConfirmation,
					[]string{loginIdentifier},
					push_email.PlaceHoldersEmailSignupConfirmation{
						ConfirmationCode: confirmationCode,
					},
				)

				if err != nil {
					return
				}
			*/
		}()
	case domainAuth.PhoneNumber:

	}

	return nil
}

func (m *Model) SignupEmailConfirmation(ctx domain.IContext, loginIdentifier, code string) (*domainAuth.TokenPair, error) {
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

func (m *Model) CreatePassword(ctx domain.IContext, email, password string) (*domainAuth.TokenPair, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tx, err := m.app.Database().Write().Begin()
	if err != nil {
		return nil, m.app.Logger().DBLog(err)
	}
	defer tx.RollbackUnlessCommitted()

	idUser, err := m.streaming.CreateUser(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err = m.repository.CreatePassword(ctx, tx, idUser, email, hashedPassword); err != nil {
		return nil, err
	}

	userDetails := &domainAuth.UserDetails{
		IdUser: idUser,
		Email:  email,
	}
	tokenPair, err := m.generateTokenPair(userDetails)

	if err = tx.Commit(); err != nil {
		return nil, m.app.Logger().DBLog(err)
	}

	return tokenPair, nil
}

func (m *Model) Logout(ctx domain.IContext, idUser int) error {
	return nil
}

func (m *Model) RefreshToken(ctx domain.IContext, refreshToken string) (*domainAuth.TokenPair, error) {
	return nil, nil
}

func (m *Model) checkIfLoginIdentifierExists(ctx domain.IContext, loginIdentifier, identifierType string) error {

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

func (m *Model) generateTokenPair(userDetails *domainAuth.UserDetails) (*domainAuth.TokenPair, error) {
	accessTokenClaims := map[string]interface{}{
		domainAuth.IdUser:    userDetails.IdUser,
		domainAuth.FirstName: userDetails.FirstName,
		domainAuth.LastName:  userDetails.LastName,
		domainAuth.Email:     userDetails.Email,
	}
	accessToken, err := m.generateTokenFromClaim(accessTokenClaims, domainAuth.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := map[string]interface{}{
		domainAuth.IdUser:    userDetails.IdUser,
		domainAuth.FirstName: userDetails.FirstName,
		domainAuth.LastName:  userDetails.LastName,
		domainAuth.Email:     userDetails.Email,
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
