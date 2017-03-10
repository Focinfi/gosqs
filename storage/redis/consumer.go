package redis

import "github.com/Focinfi/sqs/models"

// Consumer for consumer
type Consumer struct {
	C *models.Client
	P int
}

// Client returns the client
func (c *Consumer) Client() *models.Client {
	return c.C
}

// SetClient set the client
func (c *Consumer) SetClient(client *models.Client) {
	c.C = client
}

// IncPriority set the priority
func (c *Consumer) IncPriority(p int) {
	c.P += p
}

// Priority returns the priority for the c
func (c *Consumer) Priority() int {
	return c.P
}

// NewConsumer returns a new Consumer based on the client
func NewConsumer(client *models.Client, priority int) *Consumer {
	return &Consumer{
		C: client,
		P: priority,
	}
}
