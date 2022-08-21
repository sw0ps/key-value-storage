package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const port = "8080"
const mainURLPath = "/v1/{key}"
const transactPath = "logs/"
const transactFile = "transaction.log"

var logger TransactionLogger

func main() {

	err := initializaTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc(mainURLPath, keyValuePutHandler).Methods(http.MethodPut)
	r.HandleFunc(mainURLPath, keyValueGetHandler).Methods(http.MethodGet)
	r.HandleFunc(mainURLPath, keyValueDeleteHandler).Methods(http.MethodDelete)

	log.Printf("Server started at port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initializaTransactionLog() error {
	var err error

	if _, err := os.Stat(transactPath); os.IsNotExist(err) {
		err := os.Mkdir(transactPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	logger, err = NewFileTransactionLogger(transactPath + transactFile)
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	return err
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

	logger.WritePut(key, string(value))

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

	logger.WriteDelete(key)
}
