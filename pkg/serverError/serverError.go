package serverError

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ServerError struct {
	Message    string
	StatusCode int
}

func (e ServerError) Error() string {
	return e.Message
}

func ErrorResponse(w http.ResponseWriter, err error) {

	var serverError ServerError
	var httpStatusCode int
	var message string

	ok := errors.As(err, &serverError)
	if !ok {
		httpStatusCode = http.StatusInternalServerError
	} else {
		httpStatusCode = serverError.StatusCode
	}

	message = err.Error()

	w.WriteHeader(httpStatusCode)

	if message != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp := make(map[string]string)
		resp["error"] = message
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
	}

}
