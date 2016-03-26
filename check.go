package heartbeat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"sync"
)

type Check struct {
	httpCheck   []Item
	scriptCheck []Item
	clusters    map[string][]Node
}

func New() *Check {
	check := new(Check)
	check.httpCheck = []Item{}
	check.scriptCheck = []Item{}
	check.clusters = map[string][]Node{}
	return check
}

func (check *Check) AddHTTPCheck(title, url string) {
	newItem := Item{
		title:     title,
		checkType: "http",
		status:    "healthy",
		target:    url,
	}
	check.httpCheck = append(check.httpCheck, newItem)
}

func (check *Check) AddScriptCheck(title, url string) {
	newItem := Item{
		title:     title,
		checkType: "script",
		status:    "healthy",
		target:    url,
	}
	check.httpCheck = append(check.httpCheck, newItem)
}
func (check *Check) Run() {
	check.run()
}

// CheckHTTP method for checking health by http
func (check *Check) CheckHTTP() ([]byte, error) {
	for _, value := range check.httpCheck {
		resp, err := check.checkItem(value.target)
		if err != nil {
			value.status = "unhealthy"
			return []byte{}, errors.New(fmt.Sprintf("Unhealthy on %s", value.target))
		} else {
			value.status = "healthy"
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				value.status = "unhealthy"
			} else {
				value.body = contents
			}

			resp.Body.Close()
		}
	}

	return json.Marshal(check.httpCheck)
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

func (check *Check) checkCluster() error {
	return nil
}
