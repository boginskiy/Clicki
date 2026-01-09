package handler_test

import (
	"context"
	"testing"

	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/rpc"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/tester/tserv"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestShortenerService(t *testing.T) {
	pathToLogg := "test.log"
	srv := InitShortenerService(pathToLogg)

	testShortenURL(t, srv)

	defer tfunc.DeleteTestFiles(pathToLogg)

}

func testShortenURL(t *testing.T, srv *handler.ShortenerService) {
	tests := []struct {
		name    string
		bodyReq string
		codeRes int
		dataRes int
	}{
		{"test positive ShortenURL gRPC request", "https://github.com", 0, 0},
		{"test negative repeate url ShortenURL gRPC request", "https://github.com", 6, 0},
		{"test negative ShortenURL gRPC request", "host", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Context
			ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "XXX")
			ctx = context.WithValue(context.Background(), auth.CtxUserID, 100)
			// Header
			// var header metadata.MD

			res, err := srv.ShortenURL(ctx, &rpc.URLShortenRequest{Url: tt.bodyReq})

			// Check error.
			grpcStatus, ok := status.FromError(err)
			if ok {

				if int(grpcStatus.Code()) != tt.codeRes {
					t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.codeRes, int(grpcStatus.Code()))
				}

			} else {
				// Check response.
				if len(res.Result) == tt.dataRes {
					t.Errorf("%s:\n\texpected: >%v\n\tactual: %v", tt.name, tt.dataRes, len(res.Result))
				}
			}
		})
	}
}

func InitShortenerService(path string) *handler.ShortenerService {
	logg := logg.NewLogg(path, "INFO")

	config := tserv.InitConfig()
	funcer := preparation.NewFunctions(logg)
	protGRPC := protocol.NewProtocolGRPC(funcer)

	URLSrv := tserv.InitURLServ(logg, config)
	APISrv := tserv.InitAPIURLServ(logg, config)

	return handler.NewShortenerService(APISrv, URLSrv, protGRPC)
}
