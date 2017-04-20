package external

// User represents a user of sqs
type User interface {
	ID(args ...interface{}) (int64, error)
}

// UserFunc user function
type UserFunc func(args ...interface{}) (int64, error)

// ID implement User interface
func (f UserFunc) ID(args ...interface{}) (int64, error) {
	return f(args...)
}

// Root for root user
var Root = UserFunc(func(args ...interface{}) (int64, error) { return 1, nil })

// UserStore for user storage
type UserStore interface {
	GetUserIDByUniqueID(uniqueID string) (int64, error)
	CreateUserByUniqueID(uniqueID string) (int64, error)
}

// DefaultUserStore default store for user
var DefaultUserStore UserStore

func init() {
	DefaultUserStore = NewMySQL()
}
