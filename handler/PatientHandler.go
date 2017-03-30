package handler

import (
	"er-wait-time/rsc"
	"net/http"
	"time"
)

var clinicPatients map[int][]rsc.Patient = make(map[int][]rsc.Patient)

type PatientsResponse struct {
	ApiResponse
	Result []rsc.Patient `json:"result"`
}

func NewPatientsResponse(a ApiResponse, r []rsc.Patient) PatientsResponse {
	return PatientsResponse{ApiResponse: a, Result: r}
}

func GetPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		instantiatePatients(clinicId)
		patients := clinicPatients[clinicId]
		returnPatients(w, r, patients)
	}
}

func AddPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		instantiatePatients(clinicId)
		patients := clinicPatients[clinicId]
		newNumber := len(patients) + 1
		if newNumber > 1 {
			newNumber = patients[newNumber-2].PatientNumber + 1
		}
		patient := rsc.Patient{PatientNumber: newNumber, CheckInTime: time.Now()}
		clinicPatients[clinicId] = append(patients, patient)
		returnPatients(w, r, []rsc.Patient{patient})
	}
}

func returnPatients(w http.ResponseWriter, r *http.Request, cPatients []rsc.Patient) {
	patients := NewPatientsResponse(NewApiResponse(200, nil), cPatients)
	EncodeHelper(w, patients)
}

func instantiatePatients(clinicId int) {
	if clinicPatients[clinicId] == nil {
		clinicPatients[clinicId] = []rsc.Patient{}
	}
}
