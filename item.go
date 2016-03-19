package heartbeat


type Item struct {
	title  string
	checkType   string
	target string
	status string
}

func (item*Item) Status() string{
	return item.status
}

