package prioritylist

import (
	"container/heap"
	"sync"
)

var mux sync.RWMutex

// Item for item a prioritylist
type Item struct {
	index    int
	Priority int
	Value    interface{}
}

type list []*Item

// PriorityList act as a prioritylist
type PriorityList struct {
	list *list
}

// Len returns the length
func (pl PriorityList) Len() int {
	return pl.Len()
}

// Pop pops the item with the heighest priority
func (pl *PriorityList) Pop() *Item {
	return heap.Pop(pl.list).(*Item)
}

// Push pushes the item
func (pl *PriorityList) Push(item *Item) {
	heap.Push(pl.list, item)
}

// NewItem return a new Item
func NewItem(value interface{}, priority int) *Item {
	return &Item{
		Value:    value,
		Priority: priority,
	}
}

// Len returns the length oflist
func (l list) Len() int {
	mux.RLock()
	defer mux.RUnlock()
	return len(l)
}

// Less returns less resoult based on priority
func (l list) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	mux.RLock()
	defer mux.RUnlock()
	return l[i].Priority > l[j].Priority
}

// Swap swaps the Consumer indexed with i and j
func (l list) Swap(i, j int) {
	mux.Lock()
	defer mux.Unlock()

	l[i], l[j] = l[j], l[i]
	l[i].index = i
	l[j].index = j
}

// Push x into pq
func (l *list) Push(x interface{}) {
	mux.Lock()
	defer mux.Unlock()

	n := len(*l)
	item := x.(*Item)
	item.index = n
	*l = append(*l, item)
}

// Pop pops returns the highest priority Consumer
func (l *list) Pop() interface{} {
	mux.Lock()
	defer mux.Unlock()

	old := *l
	n := len(old)
	consumer := old[n-1]
	consumer.index = -1 // for safety
	*l = old[0 : n-1]
	return consumer
}
