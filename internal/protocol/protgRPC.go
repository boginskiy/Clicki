package protocol

import (
	"github.com/boginskiy/Clicki/internal/grpc"
	"github.com/boginskiy/Clicki/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProtocolGRPC struct {
}

func NewProtocolGRPC() *ProtocolGRPC {
	return &ProtocolGRPC{}
}

func (s *ProtocolGRPC) GetURLFromRequest(req any) (*model.URLJson, error) {
	in, ok := req.(*grpc.URLShortenRequest)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "bad request")
	}
	return &model.URLJson{URL: in.Url}, nil
}
