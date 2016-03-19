package heartbeat

import (
   "net/http"
)

type Handler struct {
	check Check
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	result, err := check.CheckHTTP()
	http.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}