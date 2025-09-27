package auther

import "net/http"

type Auther interface {
	NextUser() int
	CreateJWT(int) (string, error)
	CreateCookie(string, string) *http.Cookie
	GetIDAndValidJWT(string) (int, error)
}
