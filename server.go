package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = "8080"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", helloGoHandler)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello http!"))
}
