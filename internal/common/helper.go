package common

import (
	"encoding/json"
	"log"
	"net/http"
)

/* This is a helper function to encode JSON in the HTTP response
Arguments:
	response http.ResponseWriter - HTTP response writer
	httpStatus int - HTTP status
	result interface{} - Result to be encoded in JSON format for the HTTP response
*/
func encodeResponse(response http.ResponseWriter, result interface{}, statusCode int) {
	response.Header().Add("content-type", "application/json")
	response.WriteHeader(statusCode)
	if result != nil {
		err := json.NewEncoder(response).Encode(result)
		if err != nil {
			log.Println(err)
			errorMsg := []byte("Error while encoding the response")
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write(errorMsg)
			return
		}
	}
}
