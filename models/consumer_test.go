package models

import "testing"

func TestPriorityConsumer(t *testing.T) {
	queue := PriorityConsumer{}
	queue.Push(NewConsumer(&queue, &Client{}, 10))
}
