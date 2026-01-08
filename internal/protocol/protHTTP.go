package protocol

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
)

type ProtocolHTTP struct {
	Funcer prep.Funcer
}

func NewProtocolHTTP(fancer prep.Funcer) *ProtocolHTTP {
	return &ProtocolHTTP{Funcer: fancer}
}

func (s *ProtocolHTTP) GetURLFromRequest(req any) (*model.URLJson, error) {
	r, ok := req.(*http.Request)
	if !ok {
		return nil, ErrDataNotValid
	}
	urlJSON := model.NewURLJson()
	err := json.NewDecoder(r.Body).Decode(urlJSON)

	if err != nil {
		return nil, err
	}
	return urlJSON, nil
}

func (s *ProtocolHTTP) PreparResult(modURLTb *model.URLTb) []byte {
	return s.Funcer.Serialization(map[string]string{"result": modURLTb.ShortURL})
}
