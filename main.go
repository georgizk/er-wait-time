package main

import (
	"er-wait-time/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	address := "0.0.0.0:8080"
	handler := NewHttpHandler()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}))(handler)

	loggingHandler := handlers.LoggingHandler(os.Stdout, corsHandler)

	http.ListenAndServe(address, loggingHandler)
}

func NewHttpHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/clinics", handler.GetClinics()).Methods("GET")
	router.HandleFunc("/clinics", handler.AddClinic()).Methods("POST")

	router.HandleFunc("/clinics/{clinicId:[0-9]+}/wait_time", handler.GetEstimedWaitTime()).Methods("GET")

	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients", handler.GetPatients()).Methods("GET")
	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients", handler.AddPatient()).Methods("POST")
	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients/{patientNumber:[0-9]+}", handler.RemovePatient()).Methods("DELETE")

	return router
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}
}
