package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type contextKey string

const userPublicKey string = "userPublicKey"
const userRoleKey string = "userRoleKey"

func createJsonErrorResponse(w http.ResponseWriter, err error) {
	// rewrite error to avoid disclosing inside problems
	err = errors.New("can't process request")
	jsonErrorResponse(w, err, http.StatusInternalServerError)
}

func jsonErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	message := ""
	if err != nil {
		message = err.Error()
	}
	jsonResponse(w, map[string]string{"message": message}, statusCode)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, string(jsonData))
}
