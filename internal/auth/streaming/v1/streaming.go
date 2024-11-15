package v1

import (
	"github.com/google/uuid"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	ctxDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http/rabbitmq/envelop/producer"
	domainAuth "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"
)

type Streaming struct {
	app domain.IApp
}

func NewStreaming(app domain.IApp) domainAuth.IStreaming {
	return &Streaming{
		app: app,
	}
}

func (s *Streaming) CreateUser(ctx ctxDomain.IContext, uuid uuid.UUID) (int, error) {
	//Todo: be applied via grpc to user service

	s.app.Grpc().x

	return 0, nil
}

func (s *Streaming) GetUserDetails(ctx ctxDomain.IContext, idUser int) (*domainAuth.UserDetails, error) {

	return &domainAuth.UserDetails{
		IdUser:    idUser,
		FirstName: "Guilherme",
		LastName:  "Alegre",
	}, nil

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
