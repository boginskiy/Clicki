package protocol

import (
	"context"
	"strconv"

	"github.com/boginskiy/Clicki/internal/grpc"
	"github.com/boginskiy/Clicki/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func (s *ProtocolGRPC) GetUserIDFromCtx(ctx context.Context) (int, error) {
	var userID int

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get("userID")
		if len(val) > 0 {
			userID, _ = strconv.Atoi(val[0])
		}
	}

	if userID == 0 {
		return userID, status.Error(codes.DataLoss, "missing user id")
	}
	return userID, nil
}
