package main

import (
	"errors"
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
	r.HandleFunc(mainUrlPath, keyValueGetHandler).Methods(http.MethodGet)
	r.HandleFunc(mainUrlPath, keyValueDeleteHandler).Methods(http.MethodDelete)

	log.Printf("Server started at port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// keyValuePutHandler Example:
// 	curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a

// keyValuePutHandler - create key value using PUT-request
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

// keyValueGetHandler Example:
//	curl -X GET -v http://localhost:8080/v1/key-a

// keyValueGetHandler - get value by key using GET-request
func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get key from request
	key := vars["key"]

	value, err := Get(key)              // get value by key
	if errors.Is(err, ErrorNoSuchKey) { // check for 404
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // write value to response
}

// keyValueDeleteHandler Example:
//	curl -X DELETE -v http://localhost:8080/v1/key-a

// keyValueDeleteHandler - delete value by key using DELETE-request
func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get key from request
	key := vars["key"]

	err := Delete(key) // delete value by key
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
