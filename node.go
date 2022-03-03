package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// DATA
var self_node_name = ""
var self_node_id = 0
var total_nodes = 0
var node_info_map = make(map[string]node)
var node_id_to_name = make(map[int]string)

type node struct {
	node_name         string
	node_id           int
	host_name         string
	ip_addr           string
	port_num          string
	connected_to_self bool
	is_ready          bool // connected to all nodes in the network
}

// CONNECTIONS
type curr_conns_mutex struct {
	mutex      sync.Mutex
	curr_conns int
}

var curr_conns = curr_conns_mutex{curr_conns: 0}
var total_conns = 0 // Total number of connections we're expecting

// THREADING
type job_queue_mutex struct {
	mutex     *sync.Mutex
	job_queue []message
	cond      *sync.Cond
}

type counter_mutex struct {
	mutex   *sync.Mutex
	counter int
}

// key: node name, val: job queue mutex struct
var job_queues = make(map[string]job_queue_mutex)
var counter = counter_mutex{counter: 0}

// Message that is sent between processes
type message struct {
	Data string
	// DeliveredIds   []int // the process ids where the message is delivered
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

// https://pkg.go.dev/container/heap
type PriorityQueue []*heap_message

type pq_mutex struct {
	pq    PriorityQueue
	mutex sync.Mutex
}

var pq_init PriorityQueue
var pq = pq_mutex{pq: pq_init}

type proposal_mutex struct {
	mutex        sync.Mutex
	sequence_num int
}

var sequence_num = proposal_mutex{sequence_num: 0}

type message_info_mutex struct {
	mutex            sync.Mutex
	message_info_map map[string]message
}

var message_info_map = message_info_mutex{ // message_info_map: {message_id: message_info_mutex}
	message_info_map: make(map[string]message),
}

var message_id_to_heap_message = make(map[string](*heap_message))

// TRANSACTIONS
type account_mutex struct {
	account_name string
	balance      int
}

type bank_mutex struct {
	mutex sync.Mutex
	bank  map[string]int
}

var bank = bank_mutex{bank: make(map[string]int)}

/////// MAIN & SETUP ///////

func main() {
	// Argument parsing
	if len(os.Args) < 3 {
		print("Incorrect number of Arguments!\n")
		os.Exit(1)
	}

	args := os.Args[1:]
	node_name := args[0]
	config_file := args[1]
	self_node_name = node_name

	// File Parsing
	content := parse_file(config_file)

	// Node Creation
	create_node_data(content)

	// Transaction setup
	create_queues()

	// Connections
	total_conns = (total_nodes - 1) * 2
	self_node := node_info_map[node_name]

	print("begin threading")
	// Threading Begins
	// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
	wg.Add(2)
	// Servers
	go recieve_conn_reqs(self_node.port_num)

	// Clients
	go send_conn_reqs(self_node.node_name)

	// Handle transactions from generator.py
	// time.sleep(6 * time.Second)
	// go add_transactions_to_queues(self_node.node_name)

	wg.Wait()
}

func parse_file(config_file string) []string {
	content, err := os.ReadFile(config_file)
	handle_err(err)

	content2 := string(content)
	content3 := strings.Split(content2, "\n")
	print(content3)

	return content3
}

func create_node_data(content []string) {
	// Node creation
	total_nodes, _ = strconv.Atoi(content[0])

	for i := 1; i <= total_nodes; i++ {
		node_info := strings.Split(content[i], " ")

		node_name := node_info[0]
		node_id := i
		node_id_to_name[i] = node_name

		ip_addr_net, _ := net.LookupIP(node_info[1])
		ip_addr := ip_addr_net[0].String()

		new_node := node{
			node_name:         node_name,
			node_id:           node_id,
			host_name:         node_info[1],
			ip_addr:           ip_addr,
			port_num:          node_info[2],
			connected_to_self: false,
		}

		node_info_map[node_name] = new_node
	}
}

func create_queues() {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for node_name := range node_info_map {
		// https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
		job_queues[node_name] = job_queue_mutex{job_queue: []message{}, mutex: &sync.Mutex{}}
		job_queue_at_node := job_queues[node_name]
		job_queue_at_node.cond = sync.NewCond(job_queues[node_name].mutex)
		job_queues[node_name] = job_queue_at_node
	}

}

/////// 1) CONNECTIONS ///////

// Iterate through all the nodes that arent ourselves and establish a connection as client
func send_conn_reqs(self_name string) {
	for name, info := range node_info_map {
		if name != self_name {
			host := info.host_name
			port := info.port_num
			go send_req(host, port, name)
		}
	}
}

// sends a request to establish connection
// Use preexisting thread to handle new connection
func send_req(host string, port string, name string) {
	print("sending req \n")
	var conn net.Conn

	for conn == nil {
		ip := host + ":" + port
		conn, err := net.Dial("tcp", ip)

		if err != nil {
			continue
		} else {
			wait_for_connections(conn, name, false)
			break
		}
	}

}

func recieve_conn_reqs(port string) {
	print("recieive conn req\n")
	for i := 0; i < total_conns/2; i++ {
		print("in receive conn req for loop \n")
		go recieve_req(port)
	}
}

func recieve_req(port string) {
	print("recieving...\n")
	// Listen for incoming connections
	serv_port := ":" + port
	ln, err := net.Listen("tcp", serv_port)
	handle_err(err)

	// Accept connections
	conn, err := ln.Accept()
	handle_err(err)

	// Check which node we just connected to
	remote_addr := conn.RemoteAddr().(*net.TCPAddr)
	received_ip := remote_addr.IP.String()

	var curr_node_name string

	for name, info := range node_info_map {
		if received_ip == info.ip_addr {
			curr_node_name = name
		}
	}

	// Close listener
	ln.Close()

	// Wait for all connections to be established
	wait_for_connections(conn, curr_node_name, true)
}

func wait_for_connections(conn net.Conn, node_name string, receiving bool) {
	print("waiting\n")
	// easiest thing to do: keep two connections between two nodes -> one for listening, other for writing
	// Increment current number of connections
	curr_conns.mutex.Lock()
	curr_conns.curr_conns += 1
	curr_conns.mutex.Unlock()

	for curr_conns.curr_conns < total_conns {
		time.Sleep(20 * time.Millisecond)
	}

	sec := 5
	print("Found all connections. Sleeping for + " + strconv.Itoa(sec) + "seconds...\n")

	// Sleep for a few seconds - make sure all the other nodes have established connections
	time.Sleep(5 * time.Second)

	// Move to handling transactions
	if receiving {
		handle_receiving_transactions(conn, node_name)
	} else {
		handle_sending_transactions(conn, node_name)
	}
}

////// 2) TRANSACTIONS  ///////

func handle_receiving_transactions(conn net.Conn, node_name string) {
	print("handling recieving\n")
	for {
		incoming, _ := bufio.NewReader(conn).ReadString('\n')
		if incoming == "" {
			continue
		}

		fmt.Print("Message Received:", string(incoming))

		new_message := str_to_message(incoming)
		incoming_message_id := new_message.Message_id
		incoming_node_id := new_message.Origin_id
		incoming_message_proposals := new_message.Proposals

		// Put new messages into the heap and dictionary
		_, ok := message_info_map.message_info_map[incoming_message_id]

		if !ok {
			// dictionary
			message_info_map.mutex.Lock()
			message_info_map.message_info_map[incoming_message_id] = new_message
			message_info_map.mutex.Unlock()

			// Priqueue
			counter.mutex.Lock()
			h := heap_message{
				message_id: incoming_message_id,
				index:      counter.counter,
				priority:   float64(sequence_num.sequence_num) + (0.1 * float64(self_node_id)),
			}

			message_id_to_heap_message[incoming_message_id] = &h

			counter.counter++
			counter.mutex.Unlock()

			pq.mutex.Lock()
			pq.pq.Push(h)
			pq.mutex.Unlock()
		}

		old_message := message_info_map.message_info_map[incoming_message_id]

		// ISIS algo
		// If origin is ourselves (receiving a proposed priority for a message we sent)
		if incoming_node_id == self_node_id {
			old_message.Proposals = combine_arrs(old_message.Proposals, incoming_message_proposals)

			// if proposals_arr = full
			if len(old_message.Proposals) == total_nodes {
				// Determine final priority
				final_pri := max_arr(old_message.Proposals)
				old_message.Final_priority = final_pri

				multicast_msg(old_message)
			}
		} else {
			// If origin was another node

			if old_message.Final_priority == -1.0 {
				// 1) Update priority array
				sequence_num.mutex.Lock()
				proposal := float64(sequence_num.sequence_num) + (0.1 * float64(self_node_id))
				old_message.Proposals = combine_arrs(old_message.Proposals, []float64{proposal})
				sequence_num.sequence_num += 1
				sequence_num.mutex.Unlock()

				// Add to jobqueue to be sent back to the original
				incoming_node_name := node_id_to_name[incoming_node_id]
				unicast_msg(old_message, incoming_node_name)
			} else {
				// 2) Priority has been determined
				old_message.Final_priority = new_message.Final_priority
				//update message in priority queue
				pq.mutex.Lock()
				pq.pq.update(message_id_to_heap_message[old_message.Message_id],
					old_message.Message_id,
					old_message.Final_priority)
				pq.mutex.Unlock()
			}
		}

		message_info_map.message_info_map[incoming_message_id] = old_message

		// Check for delivery
		deliver_messages()
	}
}

func add_transactions_to_queues(self_name string) {
	print("handling adding transactions\n")
	// read from stdin
	for {
		reader := bufio.NewReader(os.Stdin)
		curr_transaction, _ := reader.ReadString('\n')
		if curr_transaction == "" {
			continue
		}

		message_id := strconv.Itoa(rand.Int())

		sequence_num.mutex.Lock()
		proposal := float64(sequence_num.sequence_num) + (0.1 * float64(self_node_id))
		proposals := []float64{proposal}
		sequence_num.sequence_num += 1
		sequence_num.mutex.Unlock()

		curr_message := message{
			Data:           curr_transaction,
			Origin_id:      self_node_id,
			Proposals:      proposals,
			Message_id:     message_id,
			Final_priority: -1,
		}

		// add this to every job_queue
		multicast_msg(curr_message)
	}
}

func unicast_msg(msg message, node_dest string) {
	job_queue_at_node := job_queues[node_dest]

	// Put on jobqueue
	job_queue_at_node.mutex.Lock()
	job_queue_at_node.job_queue = append(job_queues[node_dest].job_queue, msg)
	job_queues[node_dest] = job_queue_at_node
	job_queue_at_node.mutex.Unlock()

	// Signal to wake up that thread
	job_queues[node_dest].cond.Signal()
}

func multicast_msg(msg message) {
	for node_name := range node_info_map {
		if node_name != self_node_name {
			// Put on jobqueue
			job_queue_at_node := job_queues[node_name]
			job_queue_at_node.mutex.Lock()
			job_queue_at_node.job_queue = append(job_queues[node_name].job_queue, msg)
			job_queues[node_name] = job_queue_at_node

			job_queue_at_node.mutex.Unlock()
			// Signal to wake up that thread
			job_queues[node_name].cond.Signal()

		}
	}
}

func handle_sending_transactions(conn net.Conn, node_name string) {
	print("handling sending\n")
	// look into condition vars, sleep/wakeup on the condition variable
	job_queues[node_name].mutex.Lock()

	curr_job_queue := job_queues[node_name].job_queue

	for len(curr_job_queue) <= 0 {
		print("No more jobs to send at " + node_name)
		job_queues[node_name].cond.Wait()
	}

	// completing a job and popping it off the jobqueue
	curr_job := curr_job_queue[0] // message struct
	conn.Write([]byte(message_to_str(curr_job)))
	curr_job_queue = curr_job_queue[1:]

	job_queues[node_name].mutex.Unlock()
}

func deliver_messages() {
	print("delivering a message\n")
	if len(pq.pq) != 0 {
		message_id_to_deliver := pq.pq.Peek().message_id
		message_to_deliver := message_info_map.message_info_map[message_id_to_deliver]

		if message_to_deliver.Final_priority > 0 {
			// Update bank
			process_message_data(message_to_deliver)

			// Update priqueue
			pq.mutex.Lock()
			pq.pq.Pop()
			pq.mutex.Unlock()
		}
	}
}

func process_message_data(m message) {
	update_bank(m)
	print_balances()
}

func update_bank(m message) {
	data := m.Data
	info := strings.Split(data, " ")

	if info[0][:1] == "T" { // Transfer
		transfer(info)
	} else if info[0][:1] == "D" { // Deposit
		deposit(info)
	}
}

func transfer(info []string) {
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
}

func deposit(info []string) {
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

////// Helpers ///////

func combine_arrs(arr1 []float64, arr2 []float64) []float64 {
	slice := append(arr1, arr2...)

	keys := make(map[float64]bool)
	list := []float64{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func message_to_str(m message) string {
	m_json, _ := json.Marshal(m)
	m_str := string(m_json)

	return m_str
}

func str_to_message(m_str string) message {
	var m message
	err := json.Unmarshal([]byte(m_str), m)
	handle_err(err)

	return m
}

func max_arr(arr []float64) float64 {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}
func handle_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

////// PRIORITY QUEUE DEFINITION //////

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
func (pq *PriorityQueue) update(item *heap_message, message_id string, priority float64) {
	item.message_id = message_id
	item.priority = priority
	heap.Fix(pq, item.index)
}

/////// Extraneous //////

// func remove_duplicates(slice []float64) []float64 {
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

// func combine_arrs_int(arr1 []int, arr2 []int) []int {
// 	slice := append(arr1, arr2...)

// 	keys := make(map[int]bool)
// 	list := []int{}

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

/////// literal garbage over here ////////

// func recieve_conn_reqs(port string) {
// 	serv_port := ":" + port
// 	ln, err := net.Listen("tcp", serv_port)
// 	handle_err(err)

// 	print("in recieve conn reqs\n")
// 	print(total_conns)
// 	print("\n")

// 	for curr_conns.curr_conns < total_conns {
// 		conn, err := ln.Accept()
// 		handle_err(err)

// 		print("established connection in recieve conn req\n")

// 		// Start thread for connection
// 		go wait_for_connections(conn)
// 	}

// 	print("out of for loop\n")

// 	// Close the listener
// 	defer ln.Close()

// 	print("closed listener\n")

// 	wg.Wait()
// }
