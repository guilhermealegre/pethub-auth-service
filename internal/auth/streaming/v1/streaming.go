package v1

import (
	"github.com/google/uuid"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/rabbitmq/envelop/producer"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
	"github.com/guilhermealegre/pethub-user-service/api/v1/grpc/user_service_user"
)

type Streaming struct {
	app        domain.IApp
	userClient user_service_user.UserClient
}

func NewStreaming(app domain.IApp, userClient user_service_user.UserClient) domainAuth.IStreaming {
	return &Streaming{
		app:        app,
		userClient: userClient,
	}
}

func (s *Streaming) CreateUser(ctx ctxDomain.IContext, uuidUser uuid.UUID) (int, error) {

	resp, err := s.userClient.CreateUser(ctx.ToGrpc(), &user_service_user.CreateUserRequest{UUID: uuidUser[:]})
	if err != nil {
		return 0, err
	}

	return int(resp.Id), nil
}

func (s *Streaming) SendEmailSignupConfirmationCode(ctx ctxDomain.IContext, email, confirmationCode string) error {

	message := &producer.NotificationSendEmailRequest{
		To:       []string{email},
		Template: producer.EmailTemplateSignupConfirmationKey,
		PlaceHolders: map[string]any{
			producer.SignupConfirmationCode: confirmationCode,
		},
	}

	err := s.app.Rabbitmq().Produce(message, producer.NotificationExchange, producer.NotificationQueueSendEmailBidding)
	if err != nil {
		return err
	}

	return nil
}

func (s *Streaming) SendSMSSignupConfirmationCode(ctx ctxDomain.IContext, number, confirmationCode string) error {

	return nil
}
