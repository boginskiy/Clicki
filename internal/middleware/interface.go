package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type MdlwereTrustee interface {
	WithTrustedSubnet(http.HandlerFunc) http.HandlerFunc
}

type MdlwereManager interface {
	WithLogg(http.HandlerFunc) http.HandlerFunc
	WithGzip(http.HandlerFunc) http.HandlerFunc
	WithAuth(http.HandlerFunc) http.HandlerFunc
}

type MdlwereConveyor interface {
	Conveyor(http.HandlerFunc) http.HandlerFunc
}

// Middleware is extra function before buisness logic.
type Middleware interface {
	MdlwereConveyor
	MdlwereManager
	MdlwereTrustee
}

type InterceptorManager interface {
	WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// Interceptor is extra function before buisness logic for gRPC.
type Interceptor interface {
	InterceptorManager
}
