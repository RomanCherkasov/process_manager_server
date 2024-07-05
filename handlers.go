package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getProcessesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get processes handler")
	processes := processManager.ListProcesses()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processes)
}

func rerunProcessHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rerun process handler")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid process ID", http.StatusBadRequest)
		return
	}

	process, err := processManager.RestartProcess(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func startProcessHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Start process handler")
	var p Process
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	process, err := processManager.StartProcess(p.Command, p.Args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Process started")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func stopProcessHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Stop process handler")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid process ID", http.StatusBadRequest)
		return
	}

	err = processManager.StopProcess(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	var ping PingPong
	if err := json.NewDecoder(r.Body).Decode(&ping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ping.Ping != "ping" {
		http.Error(w, "Invalid ping", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	pong := PingPong{Ping: "pong"}
	json.NewEncoder(w).Encode(pong)
}
