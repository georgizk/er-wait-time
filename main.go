package main

import (
	"er-wait-time/handler"
	"er-wait-time/rsc"
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

	router.HandleFunc("/clinics", handler.GetClinics()).Methods("GET")
	router.HandleFunc("/clinics", handler.AddClinic()).Methods("POST")

	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients", handler.GetPatients()).Methods("GET")
	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients", handler.AddPatient()).Methods("POST")
	router.HandleFunc("/clinics/{clinicId:[0-9]+}/patients/{patientNumber:[0-9]+}", handler.RemovePatient()).Methods("DELETE")

	return router
}

// To be kept in memory
var clinicVisitTimes map[rsc.Clinic][]time.Time

// To be executed every time that a 'sign out' happens
func calculateAvgVisitTime(clinic rsc.Clinic, patient rsc.Patient) float32 {
	visitTimes := clinicVisitTimes[clinic]
	currentVisitTime := patient.CheckoutTime - time.Now()

	averageVisitTime := 0.0
	for visitTime := range visitTimes {
		averageVisitTime += visitTime
	}
	averageVisitTime += currentVisitTime
	averageVisitTime /= len(visitTimes + 1)
	return averageVisitTime
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}
}
