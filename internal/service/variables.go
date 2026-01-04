package service

import "errors"

const (
	LONG = 8
)

var (
	// Errors:
	ErrReadRecord     = errors.New("the record is in the queue for deletion")
	ErrDataNotValid   = errors.New("data not available or invalid")
	ErrUserIDNotValid = errors.New("userID not available or invalid")
	ErrURLNotValid    = errors.New("url not available or invalid")

	// Info:
	StoreDBIsSucces = []byte("The connection to DataBase is successfully")

	// Extra struct:
	EmptyByteSlice = []byte{}
)
