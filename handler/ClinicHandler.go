package handler

import (
	"er-wait-time/rsc"
	"errors"
	"net/http"
	"sync"
)

var allClinics map[string]rsc.Clinic = make(map[string]rsc.Clinic)
var clinicMutex sync.Mutex

type ClinicsResponse struct {
	ApiResponse
	Result map[string]rsc.Clinic `json:"result"`
}

type WaitTimeResponse struct {
	ApiResponse
	Clinic       rsc.Clinic `json:"clinic"`
	NumInQueue   int        `json:"numQueued"`
	MeanWaitTime float64    `json:"meanWaitTime"`
}

func NewClinicsResponse(a ApiResponse, r map[string]rsc.Clinic) ClinicsResponse {
	return ClinicsResponse{ApiResponse: a, Result: r}
}

func NewWaitTimeResponse(a ApiResponse, r rsc.Clinic) WaitTimeResponse {
	mean := r.CalculateAvgWaitTime()
	//	if mean == 0 {
	//		mean = 1
	//	}
	//	invMean := 1 / mean
	//	numSampled := r.NumProcessedPatients
	//	rootOfSamples := math.Sqrt(float64(numSampled))
	//
	//	if rootOfSamples == 0 {
	//		rootOfSamples = 1
	//	}

	//	lower := invMean * (1 - 1.96/rootOfSamples)
	//	upper := invMean * (1 + 1.96/rootOfSamples)
	return WaitTimeResponse{ApiResponse: a, Clinic: r, MeanWaitTime: mean, NumInQueue: len(r.QueuedPatients)}
}

func GetClinic(clinicId string) (error, rsc.Clinic) {
	clinic, ok := allClinics[clinicId]
	if ok {
		return nil, clinic
	}
	return errors.New("Index out of bounds"), rsc.Clinic{}
}

func SaveClinic(clinicId string, clinic rsc.Clinic) error {
	allClinics[clinicId] = clinic
	return nil
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
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinicId := GetStringParam(r, "clinicId")
		err, clinic := GetClinic(clinicId)
		if err != nil {
			str := err.Error()
			rsp := NewApiResponse(500, &str)
			EncodeHelper(w, rsp)
			return
		}
		response := NewWaitTimeResponse(NewApiResponse(200, nil), clinic)
		EncodeHelper(w, response)
	}
}

func returnClinics(w http.ResponseWriter, r *http.Request) {
	clinics := NewClinicsResponse(NewApiResponse(200, nil), allClinics)
	EncodeHelper(w, clinics)
}

func AddClinic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinicMutex.Lock()
		defer clinicMutex.Unlock()
		clinic := rsc.NewClinic("", "")
		err := DecodeHelper(r, &clinic)
		if err != nil {
			str := err.Error()
			errRsp := NewApiResponse(500, &str)
			clinics := NewClinicsResponse(errRsp, allClinics)
			EncodeHelper(w, clinics)
		} else {
			allClinics[clinic.UUID] = clinic
			returnClinics(w, r)
		}
	}
}
