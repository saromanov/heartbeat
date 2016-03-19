package heartbeat

import (
   "sync"
   "net/http"
   "errors"
   "encoding/json"
)

type Check struct {
	httpCheck []Item
	scriptCheck []Item
}

func New()*Check {
	check := new(Check)
	check.httpCheck := map[string]Item{}
	return check
}

func (check*Check) AddHTTPCheck(title, url string) {
	newItem := Item {
		title: title,
		checkType: "http",
		status:"healthy",
		target: url,
	}
	check.httpCheck = append(check.httpCheck, newItem)
}

func (check*Check) AddScriptCheck(title, url) {
	newItem := Item {
		title: title,
		checkType: "script",
		status:"healthy",
		target: url,
	}
	check.httpCheck[title] = url
}
func (check*Check) Run() {
	check.run()
}

func (check*Check) CheckHTTP()([]byte, error) {
	for _, value := range check.httpCheck {
		resp, err := check.checkItem(value.target)
		if err != nil {
			value.status = "unhealthy"
		}
	}

	return json.Marshal(check.httpCheck)
}

func (check*Check) run(){
	var wg sync.WaitGroup

	go func(){
		for _, value := range check.httpCheck {
			_, err := check.checkItem(value.target)
			if err != nil {

			}
			
		}
	}()
}

func (check*Check) checkItem(target string)(*http.Response, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 400 {
		return resp, errors.New("Unhealthy")
	}

	return resp, nil
}