package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/processes", getProcessesHandler).Methods("GET")
	r.HandleFunc("/processes", startProcessHandler).Methods("POST")
	r.HandleFunc("/processes/{id}", stopProcessHandler).Methods("DELETE")
	r.HandleFunc("/processes/rerun/{id}", rerunProcessHandler).Methods("POST")

	r.HandleFunc("/ping", pingHandler).Methods("POST")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// /Users/roman/Library/Android/sdk
