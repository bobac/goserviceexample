package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status struct {
	Message string `json:"message"`
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request: ", r.URL.Path, " from ", r.RemoteAddr)
	status := Status{Message: "Service is running"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
	}
}

func startHttpServer() {
	http.HandleFunc("/", handleStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
