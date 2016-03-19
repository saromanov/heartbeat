package heartbeat

import (
   "sync"
   "net/http"
)

type Check struct {
	httpCheck map[string]Item
	scriptCheck map[string]Item
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
	check.httpCheck[title] = url
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

func (check*Check) run(){
	var wg sync.WaitGroup

	go func(){
		for _, value := range check.httpCheck {
			resp, err := http.Get(value.target)
			if err != nil {

			}

			if resp.StatusCode > 400 {
				value.status = "unhealthy"
			}
		}
	}()
}