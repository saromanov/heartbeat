package heartbeat

import (
	"encoding/json"
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

func makeServer() {

}
