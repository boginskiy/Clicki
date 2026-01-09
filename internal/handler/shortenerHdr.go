package handler

import (
	"context"
	"reflect"

	p "github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/rpc"
	"github.com/boginskiy/Clicki/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerService struct {
	rpc.UnimplementedShortenerServiceServer
	APIURLServ service.APIServicer // APIURLServ is business service.
	URLServ    service.Servicer    // URLServ is business service.
	Protocol   p.Protocol          // Protocol is Protocol service.
}

func NewShortenerService(
	apiURLServ service.APIServicer,
	urlServ service.Servicer,
	prot p.Protocol) *ShortenerService {

	return &ShortenerService{
		APIURLServ: apiURLServ,
		URLServ:    urlServ,
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
	// Put in "Create" obj "Protocol" for processing "request" in "APIURLServ".
	dataByte, err := s.URLServ.Read(ctx, s.Protocol, in)

	if err == service.ErrReadRecord {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &rpc.URLExpandResponse{Result: string(dataByte)}, nil
}

func (s *ShortenerService) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*rpc.UserURLsResponse, error) {
	data, err := s.APIURLServ.ReadSet(ctx, s.Protocol)

	if err != nil || reflect.TypeOf(data) != reflect.TypeOf([]*rpc.URLData(nil)) {
		return nil, status.Error(codes.InvalidArgument, "bad request")
	}

	dataURL := data.([]*rpc.URLData)
	if len(dataURL) == 0 {
		return nil, status.Error(codes.OK, "no content")
	}

	return &rpc.UserURLsResponse{Url: dataURL}, nil
}
