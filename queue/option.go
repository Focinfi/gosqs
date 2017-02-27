package queue

// Option contains setting
type option struct {
	capacity int
	name     string
	userID   int64
}

// Config configure
type Config func(opt *option)

// SetCapacity set capacity for a queue
func SetCapacity(capactiy int) Config {
	return func(opt *option) {
		opt.capacity = capactiy
	}
}

// SetName set name for the new queue
func SetName(name string) Config {
	return func(opt *option) {
		opt.name = name
	}

}
