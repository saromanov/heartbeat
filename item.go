package heartbeat


type Item struct {
	title  string  `json:"Title"`
	checkType   string `json:"CheckType"`
	target string  `json:"Target"`
	status string  `json:"Status"`
}

func (item*Item) Status() string{
	return item.status
}

