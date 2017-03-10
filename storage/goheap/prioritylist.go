package goheap

import (
	"container/heap"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

// PriorityList for priority consumsers
type PriorityList struct {
	consumers *PriorityConsumer
}

// New returns a new PriorityConsumer
func New() *PriorityList {
	return &PriorityList{
		consumers: &PriorityConsumer{},
	}
}

// Pop returns the Consumer with the highest priority
func (p *PriorityList) Pop() (models.Consumer, error) {
	if p.consumers.Len() > 0 {
		return heap.Pop(p.consumers).(*Consumer), nil
	}

	return nil, errors.NoConsumer
}

// Push pushes the Consumer in
func (p *PriorityList) Push(consumer models.Consumer) error {
	heap.Push(p.consumers, consumer)
	return nil
}
