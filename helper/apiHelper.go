package helper

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Payload interface{} `json:"payload"`
}

func WriteErrorResponse(res http.ResponseWriter, err error) {
	apiError := ApiError{
		Code:  1,
		Error: err.Error(),
	}
	bytes, err := json.Marshal(apiError)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(bytes)
}

func WriteResponse(res http.ResponseWriter, payload interface{}) {
	response := ApiResponse{
		Code:    0,
		Payload: payload,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		WriteErrorResponse(res, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}
