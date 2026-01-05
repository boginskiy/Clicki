package server

import (
	"net"

	"github.com/boginskiy/Clicki/cmd/config"
	appGRPC "github.com/boginskiy/Clicki/internal/grpc"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"google.golang.org/grpc"
)

type ServgRPC struct {
	Cfg    config.Config
	Logg   logg.Logger
	S      *grpc.Server
	listen net.Listener
}

func NewServgRPC(
	config config.Config,
	logger logg.Logger,
	service appGRPC.ShortenerServiceServer,
	interceptor mv.Interceptor) *ServgRPC {

	listen, err := net.Listen("tcp", config.GetSrvAddr())
	if err != nil {
		logger.RaiseFatal(err, "NewServgRPC>Listen", nil)
	}

	tmpSrv := &ServgRPC{
		Cfg:  config,
		Logg: logger,
		S: grpc.NewServer(
			// AuthInterceptor is auth for gRPC.
			grpc.UnaryInterceptor(interceptor.WithAuth)),
		listen: listen,
	}

	// Registration service.
	appGRPC.RegisterShortenerServiceServer(tmpSrv.S, service)
	return tmpSrv
}

func (s *ServgRPC) Run() {
	s.Logg.RaiseFatal(s.S.Serve(s.listen), "gRPC server has not started", nil)
}
