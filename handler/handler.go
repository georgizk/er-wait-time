package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ApiResponseIf interface {
	getCode() int
	getError() error
}
type ApiResponse struct {
	Code  int     `json:"-"`
	Error *string `json:"error"`
}

func NewApiResponse(code int, e *string) ApiResponse {
	return ApiResponse{Code: code, Error: e}
}

func (r ApiResponse) getCode() int {
	return r.Code
}

func (r ApiResponse) getError() error {
	if nil != r.Error {
		return errors.New(*r.Error)
	} else {
		return nil
	}
}

func DecodeHelper(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(s)
	if err != nil {
		log.Println("[---] Decode error:", err)
		return errors.New("Decode error")
	}
	return nil
}

func EncodeHelper(w http.ResponseWriter, s ApiResponseIf) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.getCode())
	encoder.Encode(s)
}
