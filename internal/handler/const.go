package handler

const (
	// Some errors
	ErrPathAndMethod = "only POST with '/' GET with '/{id}' requests are allowed"
	ErrBodyReq       = "data not available or invalid"
	ErrNotData       = "data not found"
)
