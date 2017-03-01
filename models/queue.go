package models

// Queue contains name-message map
type Queue struct {
	UserID int64
	Name   string
}

// NewQueue returns a new queue
func NewQueue(userID int64, name string) *Queue {
	return &Queue{
		UserID: userID,
		Name:   name,
	}
}
