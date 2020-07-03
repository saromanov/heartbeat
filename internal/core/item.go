package core

type Item struct {
	id         int    `json:"id"`
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
