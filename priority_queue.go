// // https://pkg.go.dev/container/heap
// package main

// import "container/heap"

// import (
// 	"container/heap"
// )

// // An Item is something we manage in a priority queue.
// type heap_message struct {
// 	message_id string // hash
// 	index      int    // The index of the item in the heap.
// 	priority   int
// 	// The index is needed by update and is maintained by the heap.Interface methods.
// }

// // A PriorityQueue implements heap.Interface and holds Items.
// type PriorityQueue []*heap_message

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
// func main() {
// 	// Some items and their priorities.
// 	items := map[string]int{
// 		"banana": 3, "apple": 2, "pear": 4,
// 	}

// 	// Create a priority queue, put the items in it, and
// 	// establish the priority queue (heap) invariants.
// 	pq := make(PriorityQueue, len(items))
// 	i := 0
// 	for message_id, priority := range items {
// 		pq[i] = &heap_message{
// 			message_id: message_id,
// 			priority:   priority,
// 			index:      i,
// 		}
// 		i++
// 	}
// 	heap.Init(&pq)

// 	// Insert a new item and then modify its priority.
// 	item := &heap_message{
// 		message_id: strconv.Itoa(rand.Int()),
// 		priority:   1,
// 	}
// 	item2 := &heap_message{
// 		message_id: strconv.Itoa(rand.Int()),
// 		priority:   1,
// 	}
// 	heap.Push(&pq, item)
// 	heap.Push(&pq, item2)
// 	pq.update(item, item.message_id, 5)

// 	// Take the items out; they arrive in decreasing priority order.
// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(*heap_message)
// 		fmt.Printf("%.2d:%s ", item.priority, item.message_id)
// 	}
// }