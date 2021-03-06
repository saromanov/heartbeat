package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saromanov/heartbeat/api"
	"github.com/saromanov/heartbeat/internal/config"
	"github.com/saromanov/heartbeat/internal/server/model"
	log "github.com/sirupsen/logrus"
)

const apiPrefix = "/api"

// Server defines server logic
type Server struct {
	check  *api.Heartbeat
	logger *log.Logger
}

// Response provides writing of response from endpoint
type Response struct {
	URL string `json:"url"`
}

func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := s.check.Stats()
	result, err := json.Marshal(stats)
	if err != nil {
		s.logger.WithError(err).Errorf("unable to marshal json")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (s *Server) addHealthCheck(w http.ResponseWriter, r *http.Request) {
	var h model.HealthCheck

	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.Title == "" {
		http.Error(w, "title is not defined", http.StatusBadRequest)
		return
	}
	if h.URL == "" {
		http.Error(w, "url is not defined", http.StatusBadRequest)
		return
	}

	if err := s.check.AddCheck(h.Title, h.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func initLogger() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.JSONFormatter{})
	l.SetLevel(log.InfoLevel)
	return l
}

// Run starting of the server
func Run(cfg *config.Config) {
	logger := initLogger()
	if cfg == nil {
		logger.Fatalf("config is not defined")
	}
	hb := api.New()
	for _, c := range cfg.Checks {
		hb.AddCheck(c.Name, c.URL)
	}
	go hb.Run(cfg.Duration)
	s := &Server{
		check:  hb,
		logger: logger,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPrefix+"/status", s.report)
	r.HandleFunc(apiPrefix+"/checks", s.addHealthCheck).Methods("POST")
	logger.Infof("server is started to %s", cfg.Address)
	logger.Fatal(http.ListenAndServe(cfg.Address, r))
}
