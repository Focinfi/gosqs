package external

// User represents a user of sqs
type User interface {
	ID() int64
}

// UserFunc user function
type UserFunc func() int64

// ID implement User interface
func (f UserFunc) ID() int64 {
	return f()
}
