package handler

import (
	"er-wait-time/rsc"
	"net/http"
)

type PatientsResponse struct {
	ApiResponse
	Result []rsc.Patient `json:"result"`
}

func NewPatientsResponse(a ApiResponse, r []rsc.Patient) PatientsResponse {
	return PatientsResponse{ApiResponse: a, Result: r}
}

func GetPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patients := clinic.QueuedPatients
		returnPatients(w, r, patients)
	}
}

func AddPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patient := clinic.AddPatient(0)
		SaveClinic(clinicId, clinic)
		returnPatients(w, r, []rsc.Patient{patient})
	}
}

func RemovePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patientNumber := GetIntParam(r, "patientNumber")
		clinic.RemovePatient(uint64(patientNumber))
		SaveClinic(clinicId, clinic)
		returnPatients(w, r, clinic.QueuedPatients)
	}
}

func returnPatients(w http.ResponseWriter, r *http.Request, cPatients []rsc.Patient) {
	patients := NewPatientsResponse(NewApiResponse(200, nil), cPatients)
	EncodeHelper(w, patients)
}
