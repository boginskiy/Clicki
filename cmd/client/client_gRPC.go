package client

import (
	"context"
	"log"
	"testing"

	"github.com/boginskiy/Clicki/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ClientGRPC struct {
	C    rpc.ShortenerServiceClient
	Conn *grpc.ClientConn
}

func NewClientGRPC(target string) *ClientGRPC {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(clientInterceptor)
	)

	if err != nil {
		log.Fatal(err)
	}

	return &ClientGRPC{
		C:    rpc.NewShortenerServiceClient(conn),
		Conn: conn,
	}
}

func (c *ClientGRPC) Close() error {
	return c.Conn.Close()
}

var (
	TOKEN string = ""
)

func TestShortenerService(t *testing.T) {
	// ClientGRPC
	client := NewClientGRPC("localhost:8080")
	defer client.Close()

	testAuth(t, client)

	if TOKEN == "" {
		t.Errorf("%s:\n\tactual: %v", "dont't taking TOKEN", TOKEN)
		return
	}

	testShortenURL(t, client)
}

func testAuth(t *testing.T, client *ClientGRPC) {
	tests := []struct {
		name          string
		authorization string
		url           string
		statusCode    int
	}{
		{"test negative authentication for gRPC", "Nemo", "", 16},
		{"test positive authentication for gRPC", "", "", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", tt.authorization)
			var header metadata.MD

			_, err := client.C.ShortenURL(ctx, &rpc.URLShortenRequest{}, grpc.Header(&header))

			grpcStatus, ok := status.FromError(err)

			// Check error.
			if ok && int(grpcStatus.Code()) != tt.statusCode {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.statusCode, int(grpcStatus.Code()))

			}

			// Take new token.
			auth := header.Get("authorization")
			if len(auth) > 0 {
				TOKEN = auth[0]
			}

		})
	}
}

func testShortenURL(t *testing.T, client *ClientGRPC) {
	// TODO ...
}
