package helper

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
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
