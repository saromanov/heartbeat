package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/saromanov/heartbeat/internal/core"
)

// Server defines server logic
type Server struct {
	check *core.Check
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

func runHeartbeat() {
	hb := core.New()
	hb.Run(1 * time.Second)
}

// Run starting of the server
func Run() {
	hb := core.New()
	hb.Run(1 * time.Second)
	s := &Server{}
	go runHeartbeat()
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.report)
	log.Fatal(http.ListenAndServe(":8100", mux))
}
