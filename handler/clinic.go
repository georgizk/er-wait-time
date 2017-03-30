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

var allClinics []rsc.Clinic

func GetClinics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clinics := NewClinicsResponse(NewApiResponse(200, nil), allClinics)
		EncodeHelper(w, clinics)
	}
}
