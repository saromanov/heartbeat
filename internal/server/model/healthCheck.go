package model

// HealthCheck defines model for health check
type HealthCheck struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
