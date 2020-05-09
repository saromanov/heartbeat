package heartbeat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

var (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

// Check provides a basic struct for checking
type Check struct {
	// list of the http checks
	httpCheck []Item
	// dict of http checks
	httpCheckMap map[string]Item
	// list of the scipt checks
	scriptCheck []Item
	clusters    map[string][]Node
}

// New provides initialization of the project
func New() *Check {
	return &Check{
		httpCheck:    []Item{},
		scriptCheck:  []Item{},
		clusters:     map[string][]Node{},
		httpCheckMap: map[string]Item{},
	}
}

// AddHTTPCheck provides adding of HTTP check
func (check *Check) AddHTTPCheck(title, url string) {
	newItem := Item{
		title:     title,
		checkType: "http",
		status:    healthy,
		target:    url,
	}
	check.httpCheckMap[title] = newItem
	check.httpCheck = append(check.httpCheck, newItem)
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

// Run provides running of the checks
func (check *Check) Run() {
	check.run()
}

// CheckHTTP method for checking health over registered http endpoints
// Return struct of results
func (check *Check) CheckHTTP() (*HTTPReport, error) {
	items := make([]HTTPItem, len(check.httpCheck))
	for _, value := range check.httpCheck {
		resp, err := check.checkItem(value.target)
		if err != nil {
			value.status = unhealthy
			items = append(items, HTTPItem{Name: value.title, Url: value.target, Error: err.Error(), Status: "down"})
		} else {
			value.status = healthy
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				items = append(items, HTTPItem{Name: value.title, Url: value.target, StatusCode: resp.Status, Status: "down"})
				value.status = unhealthy
			} else {
				value.body = contents
				items = append(items, HTTPItem{Name: value.title, Url: value.target, StatusCode: resp.Status, Status: "up"})
			}

			resp.Body.Close()
		}
	}

	return &HTTPReport{Items: items}, nil
}

// Report provides output info to console
func (check *Check) Report() {
	items, err := check.CheckHTTP()
	if err != nil {
		log.Fatal(fmt.Errorf("%v", err))
	}

	fmt.Println("Current time: ", time.Now())
	for _, item := range items.Items {
		if item.Status == "down" {
			color.Red(fmt.Sprintf("%s - %s", item.Name, item.Url))
		} else {
			color.Green(fmt.Sprintf("%s - %s", item.Name, item.Url))
		}
	}
}

// Checking provides checking of the data
func (check *Check) Checking(interval int32) {
	for {
		time.Sleep(time.Duration(interval) * time.Minute)
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

func (check *Check) run() {
	//var wg sync.WaitGroup

	go func() {
		for _, value := range check.httpCheck {
			_, err := check.checkItem(value.target)
			if err != nil {
				// It seems, that item is unhealty

			}

		}
	}()
}

func (check *Check) checkItem(target string) (*http.Response, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 400 {
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
			return errors.New(fmt.Sprintf("Cluster %s is unhealthy. %d nodes from %d is unhealthy", title, unhealthyNodes, totalNodes))
		}
	}

	return nil
}
