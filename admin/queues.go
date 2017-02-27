package admin

import "github.com/Focinfi/sqs/queue"

// AddQueue adds a queue into root queues
func AddQueue(name string, configs ...queue.Config) error {
	if _, ok := defaultService.queues[name]; ok {
		return ErrDuplicateQueue
	}

	q := queue.New(name, configs...)
	defaultService.queues[name] = q
	return nil
}
