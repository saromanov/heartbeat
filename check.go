package heartbeat

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Check struct {
	httpCheck   []Item
	scriptCheck []Item
}

func New() *Check {
	check := new(Check)
	check.httpCheck = []Item{}
	check.scriptCheck = []Item{}
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
		}
	}

	return json.Marshal(check.httpCheck)
}

func (check *Check) run() {
	var wg sync.WaitGroup

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
