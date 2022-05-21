package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {
	http.HandleFunc("/", helloGoHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello http!"))
}
