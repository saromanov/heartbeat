package core

type Item struct {
	title      string `json:"title"`
	checkType  string `json:"checkType"`
	target     string `json:"target"`
	status     string `json:"status"`
	statusCode int    `json:"statusCode"`
	body       []byte `json:"body"`
}

func (item *Item) Status() string {
	return item.status
}
