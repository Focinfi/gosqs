package prioritylist

import "container/heap"

// NewWithStringSlice return a Interface value is string
func NewWithStringSlice(strs ...string) *StringSlice {
	ls := list(make([]*Item, len(strs)))
	for i, str := range strs {
		if str == "" {
			continue
		}

		ls[i] = &Item{
			index:    i,
			Priority: 0,
			Value:    str,
		}
	}

	heap.Init(&ls)

	return &StringSlice{list: &ls}
}

// StringSlice acts as an prioritylist
type StringSlice struct {
	*list
}

// Len returns the len of list
func (ss StringSlice) Len() int {
	return ss.list.Len()
}

// Pop pops a item value
func (ss *StringSlice) Pop() string {
	return heap.Pop(ss.list).(string)
}

// Push pushes a item with string value and priority
func (ss *StringSlice) Push(value string, priority int) {
	heap.Push(ss.list, &Item{Value: value, Priority: priority})
}
