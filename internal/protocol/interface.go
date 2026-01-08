package protocol

import (
	"github.com/boginskiy/Clicki/internal/model"
)

type Preparator interface {
	PreparResult(modURLTb *model.URLTb) []byte
}

type Protocol interface {
	GetURLFromRequest(req any) (*model.URLJson, error)
	Preparator
}
