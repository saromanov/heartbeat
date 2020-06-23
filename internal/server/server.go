package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/saromanov/heartbeat/api"
)

// Server defines server logic
type Server struct {
	check *api.Heartbeat
}

func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(Response{})
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Response provides writing of response from endpoint
type Response struct {
	URL string `json:"url"`
}

// Run starting of the server
func Run() {
	hb := api.New()
	go hb.Run(1 * time.Second)
	s := &Server{
		check: hb,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.report)
	log.Fatal(http.ListenAndServe(":8100", mux))
}
