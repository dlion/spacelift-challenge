package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spacelift-io/homework-object-storage/handlers"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/object/{id}", handlers.UploadHandler).Methods("PUT")
	r.HandleFunc("/object/{id}", handlers.GetHandler).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
