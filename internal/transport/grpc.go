package authgrpc

import (
	"context"
	"errors"
	"github.com/Vyacheslav1557/ms-auth/internal/services"
	"github.com/Vyacheslav1557/ms-auth/pkg/go/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type Auth interface {
	Login(string, string) (string, error)
	Logout(string) error
	Refresh(string) (string, error)
	Register(string, string) (string, error)
}

type serverAPI struct {
	gen.UnimplementedAuthServiceServer
	log  *slog.Logger
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth, log *slog.Logger) {
	gen.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth, log: log})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *gen.LoginRequest,
) (*gen.LoginResponse, error) {
	token, err := s.auth.Login(in.GetUsername(), in.GetPassword())
	if errors.Is(err, services.ErrBadCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	if err != nil {
		s.log.Debug(err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gen.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Refresh(
	ctx context.Context,
	in *gen.RefreshRequest,
) (*gen.RefreshResponse, error) {
	tkn, err := s.auth.Refresh(in.GetToken())
	if errors.Is(err, services.ErrBadCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	if err != nil {
		s.log.Debug(err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gen.RefreshResponse{Token: tkn}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	in *gen.LogoutRequest,
) (*emptypb.Empty, error) {
	err := s.auth.Logout(in.GetToken())
	if errors.Is(err, services.ErrBadCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	if err != nil {
		s.log.Debug(err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *gen.RegisterRequest,
) (*gen.RegisterResponse, error) {
	token, err := s.auth.Register(in.GetUsername(), in.GetPassword())
	if errors.Is(err, services.ErrBadCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}
	if err != nil {
		s.log.Debug(err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gen.RegisterResponse{Token: token}, nil
}
