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

// Root for root user
var Root = UserFunc(func() int64 { return 1 })

// TestClient for test client
var TestClient = UserFunc(func() int64 { return 1 })

// GetUserWithKey returns the userID with the params
func GetUserWithKey(accessKey string, secretKey string) (int64, error) {
	// TODO: authentication
	return 1, nil
}
