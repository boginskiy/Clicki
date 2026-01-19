package handler_test

import (
	"context"
	"strings"
	"testing"

	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/rpc"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/tester/tserv"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var setTestURLs = []string{"https://github.com", "https://translate.yandex.ru/"}

func TestShortenerService(t *testing.T) {
	pathToLogg := "test.log"
	srv := InitShortenerService(pathToLogg)

	testShortenURL(t, srv)
	testExpandURL(t, srv)
	testListUserURLs(t, srv)

	defer tfunc.DeleteTestFiles(pathToLogg)
}

func testListUserURLs(t *testing.T, srv *handler.ShortenerService) {
	ctxUser101 := context.WithValue(context.Background(), auth.CtxUserID, 101) // User 101
	ctxUser100 := context.WithValue(context.Background(), auth.CtxUserID, 100) // User 100

	// Check User101 without data
	res, err := srv.ListUserURLs(ctxUser101, &emptypb.Empty{})
	if err != nil {
		t.Errorf("%s:\n\texpected: %v\n\tactual: %v", "check nil", nil, err)
	}
	if len(res.GetUrl()) > 0 {
		t.Errorf("%s:\n\texpected: %v\n\tactual: %v", "check nil", 0, len(res.GetUrl()))
	}

	// Check User100 with data
	res, err = srv.ListUserURLs(ctxUser100, &emptypb.Empty{})
	if err != nil {
		t.Errorf("%s:\n\texpected: %v\n\tactual: %v", "check nil", nil, err)
	}

	cnt := len(setTestURLs)

	for _, urlR := range res.GetUrl() {
		for _, urlT := range setTestURLs {
			if strings.Contains(urlT, urlR.OriginalUrl) {
				cnt--
			}
		}
	}

	if cnt != 0 {
		t.Errorf("%s:\n\texpected: %+v\n\tactual: %+v", "check response with data", setTestURLs, res)
	}
}

func testExpandURL(t *testing.T, srv *handler.ShortenerService) {
	tests := []struct {
		name    string
		bodyReq string
		codeRes int
		dataRes int
	}{
		{"test negative ExpandURL gRPC request", "XXXXXXXX", 3, 0},
		{"test positive ExpandURL gRPC request", "wrs4db6j", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Context
			ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)

			res, err := srv.ExpandURL(ctx, &rpc.URLExpandRequest{Id: tt.bodyReq})

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

func testShortenURL(t *testing.T, srv *handler.ShortenerService) {
	tests := []struct {
		name    string
		bodyReq string
		codeRes int
		dataRes int
	}{
		{"test positive ShortenURL gRPC request", setTestURLs[0], 0, 0},
		{"test positive ShortenURL gRPC request", setTestURLs[1], 0, 0},
		{"test repeated url for ShortenURL gRPC request", setTestURLs[0], 6, 0},
		{"test not valid url for ShortenURL gRPC request", "github.com", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Context
			ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)

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
	//
	URLSrv, APISrv := tserv.InitURLServAndAPIURLServ(logg, config)

	return handler.NewShortenerService(APISrv, URLSrv, protGRPC)
}
