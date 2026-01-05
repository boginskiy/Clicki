package handler

import (
	"context"

	"github.com/boginskiy/Clicki/internal/grpc"
	p "github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerService struct {
	grpc.UnimplementedShortenerServiceServer
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

func (s *ShortenerService) ShortenURL(ctx context.Context, in *grpc.URLShortenRequest) (*grpc.URLShortenResponse, error) {
	// Put in "Create" obj "Protocol" for processing "request" in "APIURLServ".
	dataByte, err := s.APIURLServ.Create(ctx, s.Protocol, in)

	if err != nil && len(dataByte) == 0 {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Make response.
	response := &grpc.URLShortenResponse{
		Result: string(dataByte),
	}

	// If conflict on the server.
	if err != nil && len(dataByte) > 0 {
		return response, status.Error(codes.AlreadyExists, err.Error())
	}

	return response, nil
}

func (s *ShortenerService) ExpandURL(ctx context.Context, in *grpc.URLExpandRequest) (*grpc.URLExpandResponse, error) {
	return &grpc.URLExpandResponse{}, nil
}

func (s *ShortenerService) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*grpc.UserURLsResponse, error) {
	return &grpc.UserURLsResponse{}, nil
}
