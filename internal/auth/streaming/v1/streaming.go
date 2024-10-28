package v1

import (
	"github.com/gocraft/dbr/v2"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
	"github.com/guilhermealegre/pethub-auth-service/internal/infrastructure/database"
)

type Streaming struct {
	app domain.IApp
}

func NewStreaming(app domain.IApp) domainAuth.IStreaming {
	return &Streaming{
		app: app,
	}
}

func (s Streaming) CreateUser(ctx domain.IContext, tx *dbr.Tx) (int, error) {
	//Todo: be applied via grpc to user service
	var IdUser int
	err := tx.InsertInto(database.UserTableUser).
		Columns(
			"active",
			"password_set",
		).
		Values(
			true,
			true,
		).
		Returning("id_users").
		LoadContext(ctx, &IdUser)

	if err != nil {
		return 0, s.app.Logger().DBLog(err)
	}

	return IdUser, nil
}
