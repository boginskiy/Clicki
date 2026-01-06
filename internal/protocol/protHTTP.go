package protocol

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Clicki/internal/auth"
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

func (s *ProtocolHTTP) GetUserIDFromCtx(ctx context.Context) (int, error) {
	var userID int

	UserID, ok := ctx.Value(auth.CtxUserID).(int)

	if !ok || UserID <= 0 {
		return userID, ErrUserIDNotValid
	}
	return UserID, nil
}
