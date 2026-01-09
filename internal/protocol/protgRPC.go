package protocol

import (
	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProtocolGRPC struct {
	Funcer prep.Funcer
}

func NewProtocolGRPC(fancer prep.Funcer) *ProtocolGRPC {
	return &ProtocolGRPC{Funcer: fancer}
}

func (s *ProtocolGRPC) GetURLFromRequest(req any) (*model.URLJson, error) {
	in, ok := req.(*rpc.URLShortenRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "bad request")
	}
	return &model.URLJson{URL: in.Url}, nil
}

func (s *ProtocolGRPC) PreparResult(modURLTb *model.URLTb) []byte {
	return []byte(modURLTb.ShortURL)
}

func (s *ProtocolGRPC) GetURLID(req any) (string, error) {
	in, ok := req.(*rpc.URLExpandRequest)
	if !ok {
		return "", status.Error(codes.InvalidArgument, "bad request")
	}
	return in.GetId(), nil
}
