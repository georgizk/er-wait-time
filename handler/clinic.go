package handler

import (
	"../rsc"
	"net/http"
)

type ClinicsResponse struct {
	ApiResponse
	Result []rsc.Clinic `json:"result"`
}

func NewClinicsResponse(a ApiResponse, r []rsc.Clinic) ClinicsResponse {
	return ClinicsResponse{ApiResponse: a, Result: r}
}

var allClinics []rsc.Clinic = []rsc.Clinic{}

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
		var clinic rsc.Clinic
		err := DecodeHelper(r, &clinic)
		if err != nil {
			allClinics = append(allClinics, clinic)
			returnClinics(w, r)
		} else {
			str := err.Error()
			errRsp := NewApiResponse(500, &str)
			clinics := NewClinicsResponse(errRsp, allClinics)
			EncodeHelper(w, clinics)
		}
	}
}
