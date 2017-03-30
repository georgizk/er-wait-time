package main

import (
	"./rsc"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	address := "0.0.0.0:8080"
	handler := NewHttpHandler()
	http.ListenAndServe(address, handler)
}

func NewHttpHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler()).Methods("GET")

	return router
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}
}
