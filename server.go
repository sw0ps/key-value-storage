package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = "8080"
const mainUrlPath = "/v1/{key}"

func main() {
	r := mux.NewRouter()

	r.HandleFunc(mainUrlPath, keyValuePutHandler).Methods(http.MethodPut)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

// keyValuePutHandler Example:
// 	curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a

// keyValuePutHandler - create key/value by PUT-request
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get key from request
	key := vars["key"]

	value, err := io.ReadAll(r.Body) // request's body with value
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
