//go:build wireinject
// +build wireinject

package wire

import (
	"nunu-eth/internal/handler"
	"nunu-eth/internal/repository"
	"nunu-eth/internal/server"
	"nunu-eth/internal/service"
	"nunu-eth/pkg/app"
	"nunu-eth/pkg/jwt"
	"nunu-eth/pkg/log"
	"nunu-eth/pkg/server/http"
	"nunu-eth/pkg/sid"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewCommonHandler,
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewCommonRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewCommonService,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
