package protocol

import (
	"github.com/boginskiy/Clicki/internal/model"
)

type Preparator interface {
	PreparCreatedResult(modURLTb *model.URLTb) []byte
	PreparReadResult(modURLSet []model.ResUserURLSet) any
}

type Protocol interface {
	GetURLFromRequest(req any) (*model.URLJson, error)
	GetURLID(req any) (string, error)
	Preparator
}
