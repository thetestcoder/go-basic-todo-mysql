package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func StartServer(router *mux.Router) {
	server := &http.Server{
		Handler: router,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
