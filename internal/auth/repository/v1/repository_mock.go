package v1

import (
	"github.com/stretchr/testify/mock"
)

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

type RepositoryMock struct {
	mock.Mock
}

/*
func (m *RepositoryMock) GetCachedUser(ctx domain.IContext, idUser int) (*v1.User, error) {
	args := m.Called(ctx, idUser)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.User), args.Error(1)
}

func (m *RepositoryMock) CacheUser(ctx domain.IContext, user v1.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

*/
