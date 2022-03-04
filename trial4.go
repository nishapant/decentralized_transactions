package main

import (
	"container/heap"
)

type heap_message struct {
	message_id string // hash
	index      int    // The index of the item in the heap.
	priority   float64
	// The index is needed by update and is maintained by the heap.Interface methods.
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*heap_message)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Peek() heap_message {
	if len(*pq) <= 0 {
		return heap_message{priority: -1}
	}

	return *(*pq)[0]
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *heap_message, priority float64) {
	item.priority = priority
	heap.Fix(pq, item.index)
}
