package v1

import (
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
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

func (r *Repository) GetAuthDetails(ctx ctxDomain.IContext, tx dbr.SessionRunner, loginIdentifier, identifierType string) (userAuthDetails *domainAuth.UserAuthDetails, err error) {
	builder := tx.Select(
		"uuid_user as user_uuid",
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

func (r *Repository) SaveConfirmationCode(ctx ctxDomain.IContext, loginIdentifier, code string) error {

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

func (r *Repository) GetConfirmationCode(ctx ctxDomain.IContext, loginIdentifier string) (code string, err error) {
	code, err = r.app.Redis().Client().Get(
		ctx,
		redis.SignupConfirmationCode.Format(r.app.Config().Env, loginIdentifier)).
		Result()

	if err != nil {
		return code, err
	}

	return code, nil
}

func (r *Repository) CreatePassword(ctx ctxDomain.IContext, tx session.ITx, email string, hashedPassword []byte) (uid uuid.UUID, err error) {

	uid, err = uuid.NewUUID()
	if err != nil {
		return uid, err
	}

	_, err = tx.InsertInto(database.AuthTableAuth).
		Columns(
			"email",
			"password",
			"uuid_user",
			"active",
		).
		Values(
			email,
			hashedPassword,
			uid,
			true,
		).
		ExecContext(ctx)

	if err != nil {
		return uid, r.app.Logger().DBLog(err)
	}

	return uid, nil
}
