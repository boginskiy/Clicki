package preparation

import "net/http"

// Funcer - .
type Funcer interface {
	TakeAllBodyFromReq(req *http.Request) (string, error)
	Deserialization(req *http.Request, st any) error
	GetProtocolFromReq(req *http.Request) string
	ChangePort(host, newPort string) string
	Serialization(any) ([]byte, error)
}
