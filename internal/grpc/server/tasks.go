package grpcserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/literaen/simple_project/users/internal/config"

	grpchandler "github.com/literaen/simple_project/users/internal/grpc/handler"

	userpb "github.com/literaen/simple_project/api"
	grpcserver "github.com/literaen/simple_project/pkg/grpc/server"
)

type UserGRPCServer struct {
	server *grpcserver.Server
}

func NewUserGRPCServer(cfg *config.Config, userService *grpchandler.UserHandler) *UserGRPCServer {
	srv := grpcserver.NewServer(5 * time.Second)

	go func() {
		userpb.RegisterUserServiceServer(srv.GetServer(), userService)

		err := srv.Start(context.TODO(), fmt.Sprintf(":%s", cfg.GRPC_Port))
		if err != nil {
			log.Fatalf("error while starting grpc user server: %v", err)
		}
	}()

	return &UserGRPCServer{server: srv}
}
