package heartbeat

// Info return basic information about system
type Info struct {
	NumClusters int
	NumUnhealthy int
	NumHttpChecks int
	NumScriptChecks int
}