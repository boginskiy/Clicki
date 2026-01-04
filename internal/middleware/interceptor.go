package middleware

import (
	"context"
	"strconv"

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

	// Create new metadata.
	newMD := metadata.New(map[string]string{
		"userID": strconv.Itoa(UserID), // Add userID.
		"token":  token,                // Add token.
	})

	// Update context.
	ctx = metadata.NewOutgoingContext(ctx, newMD)
	return handler(ctx, req)
}
