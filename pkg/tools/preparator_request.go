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

func (p *PreparatorRequest) ChangePort(host, newPort string) string {
	tmpSl := strings.Split(host, ":")
	tmpSl[1] = newPort
	return strings.Join(tmpSl, "")
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
