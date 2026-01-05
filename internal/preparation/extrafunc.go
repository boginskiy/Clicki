package preparation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/internal/logg"
)

// Functions - struct with some functions.
type Functions struct {
	Logg logg.Logger
}

func NewFunctions(logger logg.Logger) *Functions {
	return &Functions{Logg: logger}
}

func (p *Functions) ChangePort(host, newPort string) string {
	tmpSl := strings.Split(host, ":")
	tmpSl[1] = newPort
	return strings.Join(tmpSl, "")
}

func (p *Functions) GetProtocolFromReq(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

func (p *Functions) TakeAllBodyFromReq(req *http.Request) (string, error) {
	originURL, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("body of request is not valid")
	}
	return string(originURL), nil
}

func (p *Functions) Deserialization(req *http.Request, st any) error {
	dec := json.NewDecoder(req.Body)
	return dec.Decode(st)
}

func (p *Functions) Serialization(st any) []byte {
	dataByte, err := json.Marshal(st)
	if err != nil {
		p.Logg.RaiseError(err, "error in serialization obj", logg.Fields{"obj": st})
		return []byte{}
	}
	return dataByte
}
