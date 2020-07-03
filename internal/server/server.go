package server

import (
	"encoding/json"
	"net/http"

	"github.com/saromanov/heartbeat/api"
	"github.com/saromanov/heartbeat/internal/config"
	log "github.com/sirupsen/logrus"
)

const apiPrefix = "/api"

// Server defines server logic
type Server struct {
	check *api.Heartbeat
}

func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(Response{})
	if err != nil {
		log.WithError(err).Errorf("unable to marshal json")
		return
	}
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
		log.Fatalf("config is not defined")
	}
	hb := api.New()
	go hb.Run(cfg.Duration)
	s := &Server{
		check: hb,
	}
	mux := http.NewServeMux()
	mux.HandleFunc(apiPrefix+"/status", s.report)
	log.Infof("server is started to %s", cfg.Address)
	log.Fatal(http.ListenAndServe(cfg.Address, mux))
}
