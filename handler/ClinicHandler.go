package handler

import (
	"er-wait-time/rsc"
	"net/http"
)

var allClinics []rsc.Clinic = []rsc.Clinic{}

type ClinicsResponse struct {
	ApiResponse
	Result []rsc.Clinic `json:"result"`
}

func NewClinicsResponse(a ApiResponse, r []rsc.Clinic) ClinicsResponse {
	return ClinicsResponse{ApiResponse: a, Result: r}
}

func GetClinics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnClinics(w, r)
	}
}

func returnClinics(w http.ResponseWriter, r *http.Request) {
	clinics := NewClinicsResponse(NewApiResponse(200, nil), allClinics)
	EncodeHelper(w, clinics)
}

func AddClinic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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