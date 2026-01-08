package handler

import (
	"context"

	p "github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/rpc"
	"github.com/boginskiy/Clicki/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerService struct {
	rpc.UnimplementedShortenerServiceServer
	APIURLServ service.ServicerTest // APIURLServ is business service.
	Protocol   p.Protocol           // Protocol is Protocol service.
}

func NewShortenerService(
	apiURLServ service.ServicerTest,
	prot p.Protocol) *ShortenerService {

	return &ShortenerService{
		APIURLServ: apiURLServ,
		Protocol:   prot,
	}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, in *rpc.URLShortenRequest) (*rpc.URLShortenResponse, error) {
	// Put in "Create" obj "Protocol" for processing "request" in "APIURLServ".
	dataByte, err := s.APIURLServ.Create(ctx, s.Protocol, in)

	if err != nil && len(dataByte) == 0 {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Make response.
	response := &rpc.URLShortenResponse{Result: string(dataByte)}

	// If conflict on the server.
	if err != nil && len(dataByte) > 0 {
		return response, status.Error(codes.AlreadyExists, err.Error())
	}

	return response, nil
}

func (s *ShortenerService) ExpandURL(ctx context.Context, in *rpc.URLExpandRequest) (*rpc.URLExpandResponse, error) {
	return &rpc.URLExpandResponse{}, nil
}

func (s *ShortenerService) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*rpc.UserURLsResponse, error) {
	return &rpc.UserURLsResponse{}, nil
}
