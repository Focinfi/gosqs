package models

// Queue contains name-message map
type Queue struct {
	UserID   int64
	Name     string
	Messages map[int64]*Message
}

// NewQueue returns a new queue
func NewQueue(userID int64, name string) *Queue {
	return &Queue{
		UserID:   userID,
		Name:     name,
		Messages: map[int64]*Message{},
	}
}
