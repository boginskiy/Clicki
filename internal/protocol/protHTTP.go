package protocol

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (s *ProtocolHTTP) PreparCreatedResult(modURLTb *model.URLTb) []byte {
	return s.Funcer.Serialization(map[string]string{"result": modURLTb.ShortURL})
}

func (s *ProtocolHTTP) GetURLID(req any) (string, error) {
	r, ok := req.(*http.Request)
	if !ok {
		return "", ErrDataNotValid
	}
	return strings.TrimLeft(r.URL.Path, "/"), nil
}

func (s *ProtocolHTTP) PreparReadResult(modURLSet []model.ResUserURLSet) any {
	if len(modURLSet) == 0 {
		return EmptyByteSlice
	}
	return s.Funcer.Serialization(modURLSet)
}
