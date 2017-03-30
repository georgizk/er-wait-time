package main

import (
	"er-wait-time/handler"
	"er-wait-time/rsc"
	"github.com/gorilla/mux"
	"net/http"
	"time"
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

var clinicVisitTimes map[rsc.Clinic][]float64 = make(map[rsc.Clinic][]float64)

func calculateAvgVisitTime(clinic rsc.Clinic, patient rsc.Patient) float64 {
	visitTimes := clinicVisitTimes[clinic]
	currentVisitTime := time.Since(patient.CheckInTime)

	var averageVisitTime float64
	for _, visitTime := range visitTimes {
		averageVisitTime += visitTime
	}
	averageVisitTime += currentVisitTime.Seconds()
	averageVisitTime /= float64(len(visitTimes) + 1)
	return averageVisitTime
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}
}
