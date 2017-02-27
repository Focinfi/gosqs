package queue

import "fmt"

// Queue contains name-message map
type Queue struct {
	*option
}

// New allocates and returns a new Queue, ready to serve
func New(configs ...Config) Queue {
	opt := &option{}
	for _, config := range configs {
		config(opt)
	}

	return Queue{
		option: opt,
	}
}

// Message contains info
type Message struct {
	UserID    int64
	QueueName string
	Content   string
}

//Name returns name
func (q Queue) Name() string {
	return q.option.name
}

// ID returns userID the q belongs to
func (q Queue) ID() string {
	return fmt.Sprintf("%d", q.option.userID)
}
