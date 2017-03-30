package handler

import (
	"er-wait-time/rsc"
	"errors"
	"net/http"
	"sync"
	"time"
)

var clinicPatients map[rsc.Clinic][]rsc.Patient = make(map[rsc.Clinic][]rsc.Patient)
var patientMutex sync.Mutex
var clinicVisitTimes map[rsc.Clinic][]float64 = make(map[rsc.Clinic][]float64)
var clinicVisitMutex sync.Mutex

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
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patientMutex.Lock()
		defer patientMutex.Unlock()
		instantiatePatients(clinic)
		patients := clinicPatients[clinic]
		returnPatients(w, r, patients)
	}
}

func AddPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patientMutex.Lock()
		defer patientMutex.Unlock()
		instantiatePatients(clinic)
		patients := clinicPatients[clinic]
		newNumber := len(patients) + 1
		if newNumber > 1 {
			newNumber = patients[newNumber-2].PatientNumber + 1
		}
		patient := rsc.Patient{PatientNumber: newNumber, CheckInTime: time.Now()}
		clinicPatients[clinic] = append(patients, patient)
		returnPatients(w, r, []rsc.Patient{patient})
	}
}

func RemovePatient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		patientNumber := GetIntParam(r, "patientNumber")
		patientMutex.Lock()
		defer patientMutex.Unlock()
		instantiatePatients(clinic)
		patients := clinicPatients[clinic]
		err, patients, removedPatient := removePatient(patients, patientNumber)
		if err == nil {
			clinicVisitMutex.Lock()
			defer clinicVisitMutex.Unlock()
			clinicPatients[clinic] = patients
			currentVisitTime := time.Since(removedPatient.CheckInTime)
			clinicVisitTimes[clinic] = append(clinicVisitTimes[clinic], currentVisitTime.Seconds())
		}
		returnPatients(w, r, clinicPatients[clinic])
	}
}

func returnPatients(w http.ResponseWriter, r *http.Request, cPatients []rsc.Patient) {
	patients := NewPatientsResponse(NewApiResponse(200, nil), cPatients)
	EncodeHelper(w, patients)
}

func instantiatePatients(clinic rsc.Clinic) {
	if clinicPatients[clinic] == nil {
		clinicPatients[clinic] = []rsc.Patient{}
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
