package tools

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

// Структура для работы с Request-запросом
type PreparatorRequest struct {
}

func (p *PreparatorRequest) ChangePort(req *http.Request, newPort string) string {
	tmpSl := strings.Split(req.Host, ":")
	tmpSl[1] = newPort
	req.Host = strings.Join(tmpSl, "")
	return req.Host
}

func (p *PreparatorRequest) CanvertBytesToString(sl []byte) string {
	return string(sl)
}

func (p *PreparatorRequest) TakeAllBodyFromReq(req *http.Request) (string, error) {
	originURL, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("body of request is not valid")
	}
	return string(originURL), nil
}

func (p *PreparatorRequest) GetProtocol(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}
