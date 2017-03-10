package goheap

import (
	"sync"

	"github.com/Focinfi/sqs/models"
)

var mux sync.RWMutex

// Consumer for a consumer pointing to a client ready to receive messages
type Consumer struct {
	client   *models.Client
	priority int
	index    int
}

// Client returns the client
func (c *Consumer) Client() *models.Client {
	return c.client
}

// SetClient set the client
func (c *Consumer) SetClient(client *models.Client) {
	c.client = client
}

// IncPriority set the priority
func (c *Consumer) IncPriority(p int) {
	c.priority += p
}

// Priority returns the priority for the c
func (c *Consumer) Priority() int {
	return c.priority
}

// NewConsumer returns a new Consumer based on the client
func NewConsumer(client *models.Client, priority int) *Consumer {
	return &Consumer{
		client:   client,
		priority: priority,
	}
}

// PriorityConsumer is a priority queue of Consumer
type PriorityConsumer []*Consumer

// Len returns the length of PriorityConsumer
func (pq PriorityConsumer) Len() int {
	mux.RLock()
	defer mux.RUnlock()
	return len(pq)
}

// Less returns less resoult based on priority
func (pq PriorityConsumer) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	mux.RLock()
	defer mux.RUnlock()
	return pq[i].priority > pq[j].priority
}

// Swap swaps the Consumer indexed with i and j
func (pq PriorityConsumer) Swap(i, j int) {
	mux.Lock()
	defer mux.Unlock()

	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push x into pq
func (pq *PriorityConsumer) Push(x interface{}) {
	mux.Lock()
	defer mux.Unlock()

	n := len(*pq)
	consumer := x.(*Consumer)
	consumer.index = n
	*pq = append(*pq, consumer)
}

// Pop pops returns the highest priority Consumer
func (pq *PriorityConsumer) Pop() interface{} {
	mux.Lock()
	defer mux.Unlock()

	old := *pq
	n := len(old)
	consumer := old[n-1]
	consumer.index = -1 // for safety
	*pq = old[0 : n-1]
	return consumer
}
