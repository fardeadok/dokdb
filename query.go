package dokdb

// QUERY FOR ASK FROM DB
// use: v := &query{.....}
type query struct {
	// q       *query
	Content string `json:"content"`
	in      chan<- object
	out     <-chan object
}

func NewQuery(c string) *query {
	v := &query{
		Content: "",
		in:      make(chan<- object),
		out:     make(<-chan object),
	}
	return v
}

// go func() {
// 	for _, n := range nums {
// 		out <- n
// 	}
// 	close(out)
// }()
