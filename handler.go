package heartbeat

import (
   "net/http"
)

type Handler struct {
	check Check
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	
}