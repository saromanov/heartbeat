package core

// Stats returns statictics for url
type Stats struct {
	URL       string `json:"url"`
	Completed uint64 `json:"completed"`
	Failed    uint64 `json:"failed"`
}
