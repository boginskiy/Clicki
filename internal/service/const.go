package service

import "errors"

const (
	DeserializFatal   = "[FATAL]: Deserialization was bad"
	SerializFatal     = "[FATAL]: Serialization was bad"
	TakeBodyReqFatal  = "[FATAL]: Taking body in request was bad"
	DataNotValidFatal = "[FATAL]: Data not available or invalid"
	LONG              = 8
)

var (
	ErrDataNotValid = errors.New("data not available or invalid")
	ConnDBIsSucces  = []byte("The connection to DataBase is successfully")
	EmptyByteSlice  = []byte{}
)
