package goheap

import (
	"container/heap"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
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

// Add adds consumer
func (p *PriorityList) Add(consumer models.Consumer) error {
	return p.Push(consumer)
}

// Push pushes the Consumer in
func (p *PriorityList) Push(consumer models.Consumer) error {
	heap.Push(p.consumers, consumer)
	log.DB.Infof("Consumers: %#v\n", p.consumers)
	return nil
}

// Remove removes the consumer
func (p *PriorityList) Remove(consumer models.Consumer) error {
	// TODO: add map into PriorityList
	return nil
}
