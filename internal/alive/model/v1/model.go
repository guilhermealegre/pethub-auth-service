package v1

import (
	"os"
	"strconv"

	v1 "github.com/guilhermealegre/pethub-auth-service/internal/alive/domain/v1"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
)

const (
	Message = "I AM ALIVE!!!"
)

type Model struct {
	app domain.IApp
}

func NewModel(app domain.IApp) v1.IModel {
	return &Model{
		app: app,
	}
}

func (m *Model) Get(ctx domain.IContext) (_ *v1.Alive, err error) {
	obj := &v1.Alive{
		ServerName: m.app.Config().Name,
		Port:       strconv.Itoa(m.app.Http().Config().Port),
		Message:    Message,
	}

	obj.Hostname, err = os.Hostname()
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (m *Model) GetPublic(ctx domain.IContext) (*v1.PublicAlive, error) {
	obj := &v1.PublicAlive{
		Name:    m.app.Config().Name,
		Message: Message,
	}

	return obj, nil
}
