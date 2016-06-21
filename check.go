package heartbeat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"sync"
)

var (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

type Check struct {
	httpCheck   []Item
	httpCheckMap map[string]Item
	scriptCheck []Item
	clusters    map[string][]Node
}

func New() *Check {
	check := new(Check)
	check.httpCheck = []Item{}
	check.scriptCheck = []Item{}
	check.clusters = map[string][]Node{}
	check.httpCheckMap = map[string]Item{}
	return check
}

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

func (check *Check) AddScriptCheck(title, url string) {
	newItem := Item{
		title:     title,
		checkType: "script",
		status:    healthy,
		target:    url,
	}
	check.httpCheck = append(check.httpCheck, newItem)
}
func (check *Check) Run() {
	check.run()
}

// CheckHTTP method for checking health over registered http endpoints
// Return struct of results
func (check *Check) CheckHTTP() ([]byte, error) {
	for _, value := range check.httpCheck {
		resp, err := check.checkItem(value.target)
		if err != nil {
			value.status = unhealthy
			return []byte{}, errors.New(fmt.Sprintf("Unhealthy on %s", value.target))
		} else {
			value.status = healthy
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				value.status = unhealthy
			} else {
				value.body = contents
			}

			resp.Body.Close()
		}
	}

	return json.Marshal(check.httpCheck)
}

// Check Cluster provides checking all clusters
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
