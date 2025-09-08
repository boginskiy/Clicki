package handler

import "net/http"

type Handler interface {
	Get(res http.ResponseWriter, req *http.Request)
	Post(res http.ResponseWriter, req *http.Request)
	Check(res http.ResponseWriter, req *http.Request)
}
