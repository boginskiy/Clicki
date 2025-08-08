package preparation

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type ExtraFuncer interface {
	TakeAllBodyFromReq(req *http.Request) (string, error)
	GetProtocolFromReq(req *http.Request) string
	ChangePort(host, newPort string) string
}

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

func (p *ExtraFunc) TakeAllBodyFromReq(req *http.Request) (string, error) {
	originURL, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("body of request is not valid")
	}
	return string(originURL), nil
}

func (p *ExtraFunc) GetProtocolFromReq(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}
