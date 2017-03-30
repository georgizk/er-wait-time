package handler

import (
	"er-wait-time/rsc"
	"errors"
	"log"
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

func RemovePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		patientNumber := GetIntParam(r, "patientNumber")
		instantiatePatients(clinicId)
		patients := clinicPatients[clinicId]
		err, patients, removedPatient := removePatient(patients, patientNumber)
		if err == nil {
			clinicPatients[clinicId] = patients
			timeDelta := time.Since(removedPatient.CheckInTime)
			log.Println("patient removed after " + timeDelta.String())
		}
		returnPatients(w, r, clinicPatients[clinicId])
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

func removePatient(s []rsc.Patient, patientNumber int) (error, []rsc.Patient, rsc.Patient) {
	for i := 0; i < len(s); i++ {
		patient := s[i]
		if patient.PatientNumber == patientNumber {

			// put removed element at the end of the array, then return
			// array that is one element shorter
			if len(s) > 1 {
				s[len(s)-1], s[i] = s[i], s[len(s)-1]
			}
			return nil, s[:len(s)-1], patient
		}
	}

	return errors.New("patient not found"), s, rsc.Patient{}
}
