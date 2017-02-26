package queue

// Queue contains name-message map
type Queue struct {
	Name     string
	Messages []Message
	*option
}

// New allocates and returns a new Queue, ready to serve
func New(name string, configs ...Config) Queue {
	opt := &option{}
	for _, config := range configs {
		config(opt)
	}

	return Queue{
		Name:     name,
		Messages: []Message{},
		option:   opt,
	}
}

// Message contains info
type Message struct {
	*Queue
	Content  interface{}
	From     string
	Recovers []string
}
