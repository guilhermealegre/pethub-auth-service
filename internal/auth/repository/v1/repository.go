package v1

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
	"github.com/guilhermealegre/pethub-auth-service/internal/infrastructure/database"
	"github.com/guilhermealegre/pethub-auth-service/internal/infrastructure/redis"
	"time"
)

type Repository struct {
	app domain.IApp
}

func NewRepository(app domain.IApp) domainAuth.IRepository {
	return &Repository{
		app: app,
	}
}

func (r *Repository) CreateAuth(ctx domain.IContext, userAuthDetails *domainAuth.UserAuthDetails) error {

	_, err := r.app.Database().Write().InsertInto(database.AuthTableAuth).
		Columns(
			"email",
			"code_phone_number",
			"phone_number",
			"email_confirmation_token",
		).
		Values(
			userAuthDetails.Email,
			userAuthDetails.CodePhoneNumber,
			userAuthDetails.PhoneNumber,
			userAuthDetails.EmailConfirmationToken,
		).
		ExecContext(ctx)

	if err != nil {
		return r.app.Logger().DBLog(err)
	}

	return nil
}

func (r *Repository) GetAuthDetails(ctx domain.IContext, tx dbr.SessionRunner, loginIdentifier, identifierType string) (userAuthDetails *domainAuth.UserAuthDetails, err error) {
	builder := tx.Select(
		"COALESCE(fk_users, 0) as id_user",
		"COALESCE(email,'') as email",
		"password",
		"COALESCE(code_phone_number, '') as code_phone_number",
		"COALESCE(phone_number,'') as phone_number",
	).
		From(database.AuthTableAuth)

	switch identifierType {
	case domainAuth.Email:
		builder.Where("email = ?", loginIdentifier)
	case domainAuth.PhoneNumber:
		builder.Where("phone_number = ?", loginIdentifier)
	}

	if _, err = builder.LoadContext(ctx, &userAuthDetails); err != nil {
		return nil, r.app.Logger().DBLog(err)
	}

	return userAuthDetails, nil
}

func (r *Repository) GetUserDetails(ctx domain.IContext, idUser int) (*domainAuth.UserDetails, error) {
	var details domainAuth.UserDetails
	column := []string{
		"id_users as id_user",
		"COALESCE(first_name, '') as first_name",
		"COALESCE(last_name,'') as last_name",
	}
	_, err := r.app.Database().Read().Select(column...).
		From(database.UserTableUser).
		Where("id_users = ?", idUser).
		LoadContext(ctx, &details)

	if err != nil {
		return nil, r.app.Logger().DBLog(err)
	}

	return &details, nil

}

func (r *Repository) SaveConfirmationCode(ctx domain.IContext, loginIdentifier, code string) error {

	_, err := r.app.Redis().Client().Set(
		ctx,
		redis.SignupConfirmationCode.Format(r.app.Config().Env, loginIdentifier),
		code,
		time.Second*60).Result()

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetConfirmationCode(ctx domain.IContext, loginIdentifier string) (code string, err error) {
	code, err = r.app.Redis().Client().Get(
		ctx,
		redis.SignupConfirmationCode.Format(r.app.Config().Env, loginIdentifier)).
		Result()

	if err != nil {
		return code, err
	}

	return code, nil
}

func (r *Repository) CreatePassword(ctx domain.IContext, tx *dbr.Tx, idUser int, email string, hashedPassword []byte) (err error) {

	_, err = tx.InsertInto(database.AuthTableAuth).
		Columns(
			"email",
			"password",
			"fk_users",
			"active",
		).
		Values(
			email,
			hashedPassword,
			idUser,
			true,
		).
		ExecContext(ctx)

	if err != nil {
		return r.app.Logger().DBLog(err)
	}

	return nil
}
