package heartbeat

import (
   "net/http"
   "encoding/json"
)

type Handler struct {
	check Check
}

func (h* Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	result, err := h.check.CheckHTTP()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}