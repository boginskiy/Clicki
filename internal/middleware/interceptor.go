package middleware

import (
	"context"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Intercept struct {
	Cfg  config.Config
	Auth auth.AutherGRPC
	Logg logg.Logger
}

func NewIntercept(config config.Config, logger logg.Logger, auther auth.AutherGRPC) *Intercept {
	return &Intercept{Cfg: config, Logg: logger, Auth: auther}
}

func (i *Intercept) WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	token, UserID, err := i.Auth.Authentication(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token is bad")
	}

	// Create header.
	header := metadata.Pairs(
		"authorization", token,
	)

	err = grpc.SetHeader(ctx, header)
	if err != nil {
		i.Logg.RaiseError(err, "error of creating headers for response", nil)
	}

	// Add UserID in the context.
	return handler(context.WithValue(ctx, auth.CtxUserID, UserID), req)
}
