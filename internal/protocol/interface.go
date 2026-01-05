package protocol

import (
	"context"

	"github.com/boginskiy/Clicki/internal/model"
)

type Protocol interface {
	GetURLFromRequest(req any) (*model.URLJson, error)
	GetUserIDFromCtx(ctx context.Context) (int, error)
}
