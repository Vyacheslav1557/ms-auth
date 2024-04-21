package app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"ms-auth/internal/services"
	authgrpc "ms-auth/internal/transport"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	log *slog.Logger,
	authService *services.Auth,
	port string,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService, log)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	var (
		err error
		l   net.Listener
	)
	l, err = net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return err
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
