package user_service_user

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type UserClientMock struct {
	mock.Mock
}

func NewUserClientMock() *UserClientMock {
	return &UserClientMock{}
}

func (c *UserClientMock) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	var optsList = []interface{}{ctx, in}
	for _, v := range opts {
		optsList = append(optsList, v)
	}

	args := c.Called(optsList...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserResponse), args.Error(1)
}
