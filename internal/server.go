package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response provides writing of response from endpoint
type Response struct {
	URL string `json:"url"`
}

func report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(Response{})
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// MakeServer provides creating of the server
func MakeServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", report)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
