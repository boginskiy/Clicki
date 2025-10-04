package auther

import "net/http"

type Auther interface {
	Authentication(req *http.Request) (*http.Cookie, int, error)
}

type JWTer interface {
	GetIDAndValidJWT(tokenStr string) (int, error)
	CreateJWT(userID int) (string, error)
}
