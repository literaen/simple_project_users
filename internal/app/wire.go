//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"time"

	"github.com/literaen/simple_project/users/internal/config"

	grpcclient "github.com/literaen/simple_project/users/internal/grpc/client"
	grpchandler "github.com/literaen/simple_project/users/internal/grpc/handler"
	grpcserver "github.com/literaen/simple_project/users/internal/grpc/server"

	"github.com/literaen/simple_project/users/internal/outbox"
	"github.com/literaen/simple_project/users/internal/users"

	"github.com/literaen/simple_project/pkg/postgres"
	"github.com/literaen/simple_project/pkg/redis"

	"github.com/google/wire"
)

type App struct {
	Config          *config.Config
	UserGRPCHandler *grpchandler.UserHandler
	UserGRPCServer  *grpcserver.UserGRPCServer
}

func InitApp() (*App, error) {
	wire.Build(
		config.LoadEnv,

		config.ProvideDBCreds,
		postgres.NewGDB,

		config.ProvideRedisCreds,
		redis.NewRDB,

		outbox.NewOutBoxRepository,
		outbox.NewOutBoxService,
		outbox.NewOutboxWorker,

		grpcclient.NewTaskGRPCClient,

		users.NewUserRepository,
		users.NewUserService,

		grpcserver.NewUserGRPCServer,
		grpchandler.NewUserHandler,

		newApp,
	)
	return nil, nil
}

func newApp(
	config *config.Config,
	gdb *postgres.GDB,
	outboxWorker *outbox.OutboxWorker,
	grpcUserHandler *grpchandler.UserHandler,
	taskGRPCServer *grpcserver.UserGRPCServer,
) *App {
	outbox.Migrate(gdb.DB)
	users.Migrate(gdb.DB)

	go outboxWorker.Start(context.TODO(), 5*time.Second, 10)

	return &App{
		Config:          config,
		UserGRPCHandler: grpcUserHandler,
		UserGRPCServer:  taskGRPCServer,
	}
}
