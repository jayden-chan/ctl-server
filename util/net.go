package util

import (
	"encoding/json"
	"net/http"
)

// HTTPStatusRes submits an HTTP response with no body
// and the specified status code
func HTTPStatusRes(res http.ResponseWriter, code int) {
	res.WriteHeader(code)
	res.Write([]byte("\n"))
}

// HTTPRes submits an HTTP resonse with the specified
// status code and body
func HTTPRes(res http.ResponseWriter, text string, code int) {
	res.WriteHeader(code)
	res.Write([]byte(text))
}

// HTTPJSONRes submits an HTTP response with JSON as the
// body and the specified status code
func HTTPJSONRes(res http.ResponseWriter, obj interface{}, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	json.NewEncoder(res).Encode(obj)
	res.Write([]byte("\n"))
}
