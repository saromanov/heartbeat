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
)

var (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

// Check provides a basic struct for checking
type Check struct {
	mu sync.RWMutex
	// list of the http checks
	httpCheck []Item
	// dict of http checks
	httpCheckMap map[string]Item
	// list of the scipt checks
	scriptCheck []Item
	stats       map[int]Stats
	clusters    map[string][]Node
}

// HTTPCheck defines check for http
type HTTPCheck struct {
	Title string
	URL   string
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
func New() *Check {
	return &Check{
		httpCheck:    []Item{},
		scriptCheck:  []Item{},
		clusters:     map[string][]Node{},
		httpCheckMap: map[string]Item{},
		stats:        map[int]Stats{},
		mu:           sync.RWMutex{},
	}
}

// AddHTTPCheck provides adding of HTTP check
func (check *Check) AddHTTPCheck(c HTTPCheck) error {
	if err := c.Validate(); err != nil {
		return err
	}
	id := len(check.httpCheck) + 1
	newItem := Item{
		id:        id,
		title:     c.Title,
		checkType: "http",
		status:    healthy,
		target:    c.URL,
	}
	check.stats[id] = Stats{
		URL: c.URL,
	}
	check.httpCheckMap[c.Title] = newItem
	check.httpCheck = append(check.httpCheck, newItem)
	return nil
}

// ApplyCheck provides applying of the check
func (check *Check) ApplyCheck(title string) error {
	item, ok := check.httpCheckMap[title]
	if !ok {
		return fmt.Errorf("Item %s is not found", title)
	}

	_, err := check.checkItem(item.target)
	if err != nil {
		return err
	}
	return nil
}

// AddScriptCheck provides adding script check
func (check *Check) AddScriptCheck(title, url string) {
	newItem := Item{
		title:     title,
		checkType: "script",
		status:    healthy,
		target:    url,
	}
	check.httpCheck = append(check.httpCheck, newItem)
}

// CheckHTTP method for checking health over registered http endpoints
// Return struct of results
func (check *Check) CheckHTTP() (*HTTPReport, error) {
	items := make([]HTTPItem, len(check.httpCheck))
	for _, value := range check.httpCheck {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		done := make(chan struct{})
		go func(id int) {
			check.mu.Lock()
			defer func() {
				check.mu.Unlock()
				done <- struct{}{}
			}()
			stats, _ := check.stats[value.id]
			resp, err := check.checkItem(value.target)
			if err != nil {
				value.status = unhealthy
				stats.Failed++
				items = append(items, HTTPItem{Name: value.title, Url: value.target, Error: err.Error(), Status: "down"})
				return
			}
			stats.Completed++
			check.stats[id] = stats
			value.status = healthy
			resp.Body.Close()
		}(value.id)

		go func(id int) {
			select {
			case <-done:
				stats, _ := check.stats[id]
				stats.Completed++
				check.stats[id] = stats
			case <-ctx.Done():
				fmt.Println("DONE")
				return
			}
		}(value.id)
	}

	return &HTTPReport{Items: items}, nil
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
	return check.stats
}

// Run provides checking
func (check *Check) Run(d time.Duration) {
	for {
		time.Sleep(d)
		check.Report()
	}
}

// CheckClusters provides checking all clusters
func (check *Check) CheckClusters() error {
	return check.checkClusters()
}

// Info return information about current checks
func (check *Check) Info() *Info {
	return &Info{
		NumClusters:   len(check.clusters),
		NumHttpChecks: len(check.httpCheck),
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
