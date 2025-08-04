package router

import (
	"net/http"

	"github.com/boginskiy/Clicki/internal/handler"
)

func ServeMux() *http.ServeMux {
	rootHandler := handler.NewRootHandler()
	mux := http.NewServeMux()
	mux.Handle("/", rootHandler)
	return mux
}
