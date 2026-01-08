package main

import (
	"context"
	"fmt"
	"log"

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

func SendRequests(ctx context.Context, client *ClientGRPC) {
	// ShortenURL
	req := &rpc.URLShortenRequest{
		Url: "https://practicum.yandex.ru1",
	}

	var header metadata.MD

	res, err := client.C.ShortenURL(ctx, req, grpc.Header(&header))

	if err != nil {
		grpcStatus, _ := status.FromError(err)
		log.Printf("bad request for ShortenURL. Mess: %v\n", grpcStatus.Message())
		return
	}

	fmt.Println(res.Result)

	fmt.Println(header.Get("authorization"))

}

func main() {
	client := NewClientGRPC("localhost:8080")
	defer client.Close()

	// Context
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njc4ODcxNjQsIlVzZXJJRCI6MX0._eynwVVtoedNZR_DYgNM-ytcgeWIqFD-EIlcSZ9aXu4",
	)

	SendRequests(ctx, client)
}
