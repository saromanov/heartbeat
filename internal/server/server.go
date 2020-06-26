package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/saromanov/heartbeat/api"
	"github.com/saromanov/heartbeat/internal/config"
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
func Run(cfg *config.Config) {
	if cfg == nil {
		panic("config is not defined")
	}
	hb := api.New()
	go hb.Run(cfg.Duration)
	s := &Server{
		check: hb,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.report)
	log.Fatal(http.ListenAndServe(":8100", mux))
}
