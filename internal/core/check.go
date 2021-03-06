package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/saromanov/heartbeat/internal/core/writer"
)

var (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

// Check provides a basic struct for checking
type Check struct {
	mu sync.RWMutex
	// list of the http checks
	httpChecks []HTTPCheck
	// dict of http checks
	httpCheckMap map[string]Item
	// list of the scipt checks
	scriptCheck []Item
	stats       map[int]Stats
	clusters    map[string][]Node
	writer      writer.Writer
	interval    time.Duration
}

// HTTPCheck defines check for http
type HTTPCheck struct {
	Title  string
	URL    string
	id     int
	status string
}

// Validate provides validating of request
func (h HTTPCheck) Validate() error {
	if h.Title == "" {
		return fmt.Errorf("title is not defined")
	}
	if h.URL == "" {
		return fmt.Errorf("url is not defined")
	}
	return nil
}

// New provides initialization of the project
func New(w writer.Writer) *Check {
	return &Check{
		httpChecks:   []HTTPCheck{},
		scriptCheck:  []Item{},
		clusters:     map[string][]Node{},
		httpCheckMap: map[string]Item{},
		stats:        map[int]Stats{},
		mu:           sync.RWMutex{},
		writer:       w,
		interval:     5 * time.Second,
	}
}

// AddHTTPCheck provides adding of HTTP check
func (check *Check) AddHTTPCheck(c HTTPCheck) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("unable to add http check: %v", err)
	}
	check.addHTTPCheck(c)
	return nil
}

// adding http check and init stats for check
func (check *Check) addHTTPCheck(c HTTPCheck) {
	c.id = len(check.httpChecks) + 1
	check.httpChecks = append(check.httpChecks, c)
	check.stats[c.id] = Stats{
		ID:    fmt.Sprintf("%d", c.id),
		URL:   c.URL,
		Title: c.Title,
	}
}

// ApplyCheck provides applying of the check
func (check *Check) ApplyCheck(title string) error {
	item, ok := check.httpCheckMap[title]
	if !ok {
		return fmt.Errorf("item %s is not found", title)
	}

	if _, err := check.checkItem(item.target); err != nil {
		return fmt.Errorf("unable to check item: %v", err)
	}
	return nil
}

// CheckHTTP method for checking health over registered http endpoints
// Return struct of results
func (check *Check) CheckHTTP() (*HTTPReport, error) {
	for _, value := range check.httpChecks {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		done := make(chan struct{})
		go func(chk HTTPCheck) {
			check.mu.Lock()
			defer func() {
				check.mu.Unlock()
				done <- struct{}{}
			}()
			resp, err := check.checkItem(value.URL)
			if err != nil {
				value.status = unhealthy
				return
			}
			value.status = healthy
			resp.Body.Close()
		}(value)

		go func(id int) {
			select {
			case <-done:
				check.IncCompleted(id)
			case <-ctx.Done():
				check.IncFailed(id)
				return
			}
		}(value.id)
	}

	return &HTTPReport{}, nil
}

// Report provides output info to console
func (check *Check) Report() {
	items, err := check.CheckHTTP()
	if err != nil {
		log.Fatal(fmt.Errorf("%v", err))
	}

	color.Red("Current time %s", time.Now().Format(time.RFC3339))
	for _, item := range items.Items {
		if item.Status == "down" {
			color.Red("%s - %s", item.Name, item.Url)
		} else {
			color.Green("%s - %s", item.Name, item.Url)
		}
	}
}

// Stats returns statistics for all endpoints
func (check *Check) Stats() map[int]Stats {
	check.mu.RLock()
	defer check.mu.RUnlock()
	return check.stats
}

// IncCompleted provides increasing of completed requests
func (check *Check) IncCompleted(id int) {
	check.mu.Lock()
	defer check.mu.Unlock()
	stats, _ := check.stats[id]
	stats.Completed++
	check.stats[id] = stats
}

// IncFailed provides increasing of failed requests
func (check *Check) IncFailed(id int) {
	check.mu.Lock()
	defer check.mu.Unlock()
	stats, _ := check.stats[id]
	stats.Failed++
	check.stats[id] = stats
}

// Run provides checking
func (check *Check) Run(d time.Duration) {
	ticker := time.NewTicker(check.interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				check.Report()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// CheckClusters provides checking all clusters
func (check *Check) CheckClusters() error {
	return check.checkClusters()
}

// Info return information about current checks
func (check *Check) Info() *Info {
	return &Info{
		NumClusters:   len(check.clusters),
		NumHttpChecks: len(check.httpChecks),
	}
}

// AddCluster provides
func (check *Check) AddCluster(name string, nodes []Node) {
	check.clusters[name] = nodes
}

func (check *Check) checkItem(target string) (*http.Response, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return resp, errors.New("Unhealthy")
	}

	return resp, nil
}

func (check *Check) checkClusters() error {
	totalNodes := len(check.clusters)
	for title, nodes := range check.clusters {
		unhealthyNodes := 0
		for _, node := range nodes {
			_, err := check.checkItem(node.Url)
			if err != nil {
				unhealthyNodes++
			}
		}

		if unhealthyNodes != 0 {
			return fmt.Errorf("Cluster %s is unhealthy. %d nodes from %d is unhealthy", title, unhealthyNodes, totalNodes)
		}
	}

	return nil
}
