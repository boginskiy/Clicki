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
	ErrReadRecord     = errors.New("the record is in the queue for deletion")
	ErrDataNotValid   = errors.New("data not available or invalid")
	ErrUserIDNotValid = errors.New("userID not available or invalid")
	StoreDBIsSucces   = []byte("The connection to DataBase is successfully")
	EmptyByteSlice    = []byte{}
)
