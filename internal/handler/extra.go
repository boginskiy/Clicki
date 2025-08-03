package handler

// Some errors
const (
	ErrPathAndMethod = "only POST with '/' GET with '/{id}' requests are allowed"
	ErrBodyReq       = "data not available or invalid"
	ErrNotData       = "data not found"
)

// Some regular expressions
const (
	CheckDomain = `^(https?:)?\/\/(www\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(\/.*)?`
	CheckPath   = `^/[a-zA-Z0-9]+$`
)

// Var
var Symbols = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
