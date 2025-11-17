package audit

import "errors"

var (
	// Регулярное выражение для определения Хендлера "GET /{id}"
	pattern = "^/[a-zA-Z0-9]{8}$"
)

var (
	ErrReadJSONBody = errors.New(`{"error":"request body has not been read"}`)
)
