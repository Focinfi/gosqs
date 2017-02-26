package queue

// Option contains setting
type option struct {
	capacity int
}

// Config configure
type Config func(opt *option)

// SetCapacity set capacity for a queue
func SetCapacity(capactiy int) Config {
	return func(otp *option) {
		otp.capacity = capactiy
	}
}
