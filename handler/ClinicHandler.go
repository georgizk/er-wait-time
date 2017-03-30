package handler

import (
	"er-wait-time/rsc"
	"errors"
	"net/http"
	"sync"
)

var allClinics []rsc.Clinic = []rsc.Clinic{}
var clinicMutex sync.Mutex

type ClinicsResponse struct {
	ApiResponse
	Result []rsc.Clinic `json:"result"`
}

type WaitTimeResponse struct {
	ApiResponse
	Clinic   rsc.Clinic `json:"clinic"`
	WaitTime float64    `json:"waitTime"`
}

func NewClinicsResponse(a ApiResponse, r []rsc.Clinic) ClinicsResponse {
	return ClinicsResponse{ApiResponse: a, Result: r}
}

func NewWaitTimeResponse(a ApiResponse, r rsc.Clinic, z float64) WaitTimeResponse {
	return WaitTimeResponse{ApiResponse: a, Clinic: r, WaitTime: z}
}

func GetClinic(clinicId int) (error, rsc.Clinic) {
	clinicMutex.Lock()
	defer clinicMutex.Unlock()
	if clinicId < len(allClinics) && clinicId >= 0 {
		return nil, allClinics[clinicId]
	}

	return errors.New("Index out of bounds"), rsc.Clinic{}
}

func GetClinics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		returnClinics(w, r)
	}
}

func GetEstimedWaitTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicId := GetIntParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		waitTime := calculateAvgWaitTime(clinic)
		response := NewWaitTimeResponse(NewApiResponse(200, nil), clinic, waitTime)
		EncodeHelper(w, response)
	}
}

func calculateAvgWaitTime(clinic rsc.Clinic) float64 {
	visitTimes := clinicVisitTimes[clinic]
	if visitTimes == nil || len(visitTimes) == 0 {
		return 0
	}

	var averageVisitTime float64
	for _, visitTime := range visitTimes {
		averageVisitTime += visitTime
	}
	averageVisitTime /= float64(len(visitTimes))

	return averageVisitTime
}

func returnClinics(w http.ResponseWriter, r *http.Request) {
	clinics := NewClinicsResponse(NewApiResponse(200, nil), allClinics)
	EncodeHelper(w, clinics)
}

func AddClinic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinic := rsc.Clinic{}
		err := DecodeHelper(r, &clinic)
		if err != nil {
			str := err.Error()
			errRsp := NewApiResponse(500, &str)
			clinics := NewClinicsResponse(errRsp, allClinics)
			EncodeHelper(w, clinics)
		} else {
			allClinics = append(allClinics, clinic)
			returnClinics(w, r)
		}
	}
}
