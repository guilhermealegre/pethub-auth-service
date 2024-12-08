package main

import (
	"fmt"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/grpc"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/rabbitmq"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/tracer"
	v1AuthController "github.com/guilhermealegre/pethub-auth-service/internal/auth/controller/v1"
	"github.com/guilhermealegre/pethub-auth-service/internal/infrastructure/providers"
	"github.com/guilhermealegre/pethub-user-service/api/v1/grpc/user_service_user"
	"os"

	v1 "github.com/guilhermealegre/pethub-auth-service/internal/auth/domain/v1"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/validator"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/database"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/app"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/http"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis"

	v1AliveController "github.com/guilhermealegre/pethub-auth-service/internal/alive/controller/v1"
	v1AliveModel "github.com/guilhermealegre/pethub-auth-service/internal/alive/model/v1"
	v1AuthModel "github.com/guilhermealegre/pethub-auth-service/internal/auth/model/v1"
	v1AuthStreaming "github.com/guilhermealegre/pethub-auth-service/internal/auth/streaming/v1"

	v1AuthRepository "github.com/guilhermealegre/pethub-auth-service/internal/auth/repository/v1"

	grpcInfra "github.com/guilhermealegre/pethub-auth-service/internal/infrastructure/grpc"
	v1Middleware "github.com/guilhermealegre/pethub-auth-service/internal/middleware/v1"
	v1SwaggerController "github.com/guilhermealegre/pethub-auth-service/internal/swagger/controller/v1"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	// app initialization

	newApp := app.New(nil)
	newHttp := http.New(newApp, nil)
	newTracer := tracer.New(newApp, nil)
	newLogger := logger.New(newApp, nil)
	newValidator := validator.New(newApp).
		AddFieldValidators().
		AddStructValidators()
	newRedis := redis.New(newApp, nil).WithAdditionalConfigType(&v1.AdditionalConfigType{})
	newDatabase := database.New(newApp, nil)
	newRabbitMQ := rabbitmq.New(newApp, nil)
	newGrpc := grpc.New(newApp, nil)

	userClient := newGrpc.GetClient(grpcInfra.UserClient)

	// repository
	authRepository := v1AuthRepository.NewRepository(newApp)

	// streaming
	authStreaming := v1AuthStreaming.NewStreaming(newApp, user_service_user.NewUserClient(userClient))

	// models
	aliveModel := v1AliveModel.NewModel(newApp)
	authModel := v1AuthModel.NewModel(newApp, authRepository, authStreaming)

	newHttp.
		//middlewares
		WithMiddleware(v1Middleware.NewAuthenticateMiddleware(newApp)).
		WithMiddleware(v1Middleware.NewPrintRequestMiddleware(newApp)).
		WithMiddleware(v1Middleware.NewPrepareContextMiddleware(newApp)).
		//controllers
		WithController(v1SwaggerController.NewController(newApp)).
		WithController(v1AliveController.NewController(newApp, aliveModel)).
		WithController(v1AuthController.NewController(newApp, authModel))

	newApp.
		WithValidator(newValidator).
		WithDatabase(newDatabase).
		WithRedis(newRedis).
		WithLogger(newLogger).
		WithTracer(newTracer).
		WithRabbitmq(newRabbitMQ).
		WithGrpc(newGrpc).
		WithHttp(newHttp)

	err := providers.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start app
	if err := newApp.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
