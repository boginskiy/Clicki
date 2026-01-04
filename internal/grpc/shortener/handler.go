package shortener

import (
	"context"

	"github.com/boginskiy/Clicki/internal/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerService struct {
	UnimplementedShortenerServiceServer

	Service grpc.ServiceProtocol
}

func NewShortenerService(serv grpc.ServiceProtocol) *ShortenerService {
	return &ShortenerService{
		Service: serv,
	}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, in *URLShortenRequest) (*URLShortenResponse, error) {
	// dataByte, err := s.APIURLServ.ShortenURL(in)

	// s.APIURLServ.Create(in)

	response := &URLShortenResponse{
		Result: "Hello, World!",
	}

	return response, nil
}

func (s *ShortenerService) ExpandURL(ctx context.Context, in *URLExpandRequest) (*URLExpandResponse, error) {
	return &URLExpandResponse{}, nil
}

func (s *ShortenerService) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*UserURLsResponse, error) {
	return &UserURLsResponse{}, nil
}
