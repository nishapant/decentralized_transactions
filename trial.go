package main

import (
	"container/heap"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// func main() {
// 	addr, err := net.LookupIP("sp22-cs425-g70-01.cs.illinois.edu")
// 	if err != nil {
// 		fmt.Println("Unknown host")
// 	} else {
// 		fmt.Println("IP address: ", addr)
// 	}
// }
// type message struct {
// 	Data           string
// 	DeliveredIds   []int // the process ids where the message is delivered
// 	OriginId       int
// 	Proposals      []float64 // null or data
// 	Message_id     string    // hash
// 	Final_priority float64   // null when start
// }

// func message_to_str(m message) string {
// 	m_json, _ := json.Marshal(m)
// 	print(m_json)
// 	m_str := string(m_json)
// 	print(m_str)
// 	return m_str
// }
// func removeDuplicateValues(slice []float64) []float64 {
// 	keys := make(map[float64]bool)
// 	list := []float64{}

// 	// If the key(values of the slice) is not equal
// 	// to the already present value in new slice (list)
// 	// then we append it. else we jump on another element.
// 	for _, entry := range slice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }
type PriorityQueue []*heap_message

type message struct {
	Data           string
	DeliveredIds   []int // the process ids where the message is delivered
	Origin_id      int
	Proposals      []float64 // null or data
	Message_id     string    // hash
	Final_priority float64   // null when start
}

type heap_message struct {
	message_id string // hash
	index      int    // The index of the item in the heap.
	priority   float64
	// The index is needed by update and is maintained by the heap.Interface methods.
}

// type account_mutex struct {
// 	account_name string
// 	balance      int
// }

type bank_mutex struct {
	mutex sync.Mutex
	bank  map[string]int
}

var bank = bank_mutex{bank: make(map[string]int)}

func main() {
	// mess := message{Data: "afeawef", DeliveredIds: []int{1, 2, 4}, OriginId: 1, Proposals: []float64{1.1, 2.2, 3.3}, Message_id: "awefawef", Final_priority: -1.0}
	// // print(mess)
	// m_json, _ := json.Marshal(mess)
	// print(m_json)
	// m_str := string(m_json)
	// print(m_str)
	// print(message_to_str(mess))
	// strings.Join(stringArray," ")
	// arr := removeDuplicateValues(append([]float64{1.2, 2.2, 3.1}, []float64{2.2, 5.4, 6.2}...))
	// print(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), " "), "[]"))
	// // print(strings.Join(arr, " "))

	// pq testing
	// Some items and their priorities.
	// items := map[string]float64{
	// 	"m1": 1, "m2": 0, "m3": 5,
	// }

	// // Create a priority queue, put the items in it, and
	// // establish the priority queue (heap) invariants.
	// pq := make(PriorityQueue, len(items))
	// i := 0
	// for value, priority := range items {
	// 	pq[i] = &heap_message{
	// 		message_id: value,
	// 		priority:   priority,
	// 		index:      i,
	// 	}
	// 	i++
	// }
	// heap.Init(&pq)

	// // Insert a new item and then modify its priority.
	// item := &heap_message{
	// 	message_id: "m4",
	// 	priority:   3,
	// }
	// heap.Push(&pq, item)
	// pq.update(item, item.message_id, 5)

	// print(pq.Peek().priority)

	// // Take the items out; they arrive in decreasing priority order.
	// for pq.Len() > 0 {
	// 	item := heap.Pop(&pq).(*heap_message)
	// 	fmt.Printf("%.2d:%s ", item.priority, item.message_id)
	// }

	// // message parse testing
	// m := message{Data: "DEPOSIT third 23", Final_priority: 2.3}
	// process_message_data(m)
	// m = message{Data: "DEPOSIT first 22", Final_priority: 2.3}
	// process_message_data(m)
	// m = message{Data: "TRANSFER first -> second 22", Final_priority: 2.3}
	// process_message_data(m)

}

func process_message_data(m message) {
	update_bank(m)
	print_balances()
}

func update_bank(m message) {
	data := m.Data
	info := strings.Split(data, " ")

	if info[0][:1] == "T" { // Transfer
		// Preprocess
		from := info[1]
		to := info[3]
		amount, _ := strconv.Atoi(info[4])

		bank.mutex.Lock()

		_, from_ok := bank.bank[from]

		// Transaction can go through
		if from_ok && bank.bank[from] >= amount {
			_, to_ok := bank.bank[to]
			if !to_ok {
				bank.bank[to] = 0
			}

			bank.bank[from] -= amount
			bank.bank[to] += amount
		}

		bank.mutex.Unlock()
	} else if info[0][:1] == "D" { // Deposit
		account := info[1]
		amount, _ := strconv.Atoi(info[2])

		bank.mutex.Lock()
		_, account_ok := bank.bank[account]
		if !account_ok {
			bank.bank[account] = 0
		}

		bank.bank[account] += amount
		bank.mutex.Unlock()
	}
}

func print_balances() {
	bank.mutex.Lock()
	balances := "BALANCES"
	accs := make([]string, 0, len(bank.bank))
	for k := range bank.bank {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		if bank.bank[acc] != 0 {
			balances += " "
			balances += acc
			balances += ": "
			amount := strconv.Itoa(bank.bank[acc])
			balances += amount
		}
	}

	balances += "\n"
	bank.mutex.Unlock()

	fmt.Println(balances)
}

//////////////////

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

func (pq *PriorityQueue) Peek() *heap_message {
	if len(*pq) <= 0 {
		return nil
	}

	return (*pq)[0]
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *heap_message, message_id string, priority float64) {
	item.message_id = message_id
	item.priority = priority
	heap.Fix(pq, item.index)
}
