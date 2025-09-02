package preparation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type ExtraFunc struct {
}

func NewExtraFunc() *ExtraFunc {
	return &ExtraFunc{}
}

func (p *ExtraFunc) ChangePort(host, newPort string) string {
	tmpSl := strings.Split(host, ":")
	tmpSl[1] = newPort
	return strings.Join(tmpSl, "")
}

func (p *ExtraFunc) GetProtocolFromReq(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

func (p *ExtraFunc) TakeAllBodyFromReq(req *http.Request) (string, error) {
	originURL, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("body of request is not valid")
	}
	return string(originURL), nil
}

func (p *ExtraFunc) Deserialization(req *http.Request, st any) error {
	dec := json.NewDecoder(req.Body)
	return dec.Decode(st)
}

func (p *ExtraFunc) Serialization(st any) ([]byte, error) {
	return json.Marshal(st)
}
