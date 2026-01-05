package protocol

import "github.com/boginskiy/Clicki/internal/model"

type Protocol interface {
	GetURLFromRequest(req any) (*model.URLJson, error)
}
