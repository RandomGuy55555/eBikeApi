package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JsonResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

/* helper function to read the body of a request and parse it into a given struct */
func ReadRequestBody(req io.ReadCloser, dataInterface interface{}) error {
	body, readErr := readRequest(req)
	if readErr != nil {
		readErrorMessage := fmt.Errorf("could not read request. %v", readErr)
		return readErrorMessage
	}

	unmarshalError := json.Unmarshal(body, dataInterface)
	if unmarshalError != nil {
		unmarshalErrorMessage := fmt.Errorf("could not unmarshal request body into given struct. %v", unmarshalError)
		return unmarshalErrorMessage
	}
	return nil
}

// reads a request and returns the body
func readRequest(request io.ReadCloser) ([]byte, error) {
	body, readRequestError := io.ReadAll(request)
	if readRequestError != nil {
		return nil, readRequestError
	}
	return body, nil
}

/*
	 function to return an error in JSON format
		1st param: the http Reponse writer
		2nd param: the error we want to return
		3rd param: the httpStatuscode we want to return
*/
func JSONError(w http.ResponseWriter, err error, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // only for dev purposes
	w.WriteHeader(httpStatusCode)
	errorResponse := JsonResponse{
		Type:    FAIL,
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(errorResponse)
}

/*
	 function to return a response in JSON format
		1st param: the http Reponse writer
		2nd param: the message we want to return
*/
func JsonSuccessResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*") // only for dev purposes
	w.WriteHeader(http.StatusOK)                       // set statuscode to 200
	response := JsonResponse{
		Type:    SUCCESS,
		Message: msg,
	}
	json.NewEncoder(w).Encode(response)
}
