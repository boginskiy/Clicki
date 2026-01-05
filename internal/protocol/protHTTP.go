package protocol

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Clicki/internal/model"
)

type ProtocolHTTP struct {
}

func NewProtocolHTTP() *ProtocolHTTP {
	return &ProtocolHTTP{}
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
