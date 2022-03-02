package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"net"
	"os"
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

// TRANSACTIONS
// Message that is sent between processes
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
	priority   int
	// The index is needed by update and is maintained by the heap.Interface methods.
}

// https://pkg.go.dev/container/heap
type PriorityQueue []*heap_message

// var pq = PriorityQueue

type proposal_mutex struct {
	mutex         sync.Mutex
	curr_proposal int
}

var curr_proposal = proposal_mutex{curr_proposal: 0}

type message_info_mutex struct {
	mutex            sync.Mutex
	message_info_map map[string]message
}

var message_info_map = message_info_mutex{ // message_info_map: {message_id: message_info_mutex}
	message_info_map: make(map[string]message),
}

// THREADING
type job_queue_mutex struct {
	mutex     *sync.Mutex
	job_queue []message
	cond      *sync.Cond
}

var job_queues = make(map[string]job_queue_mutex)

/////// MAIN ///////

func main() {
	// Argument parsing
	if len(os.Args) < 3 {
		print("Incorrect number of Arguments!\n")
		os.Exit(1)
	}

	args := os.Args[1:]
	node_name := args[0]
	config_file := args[1]

	print("finished arg parsing\n")

	// File Parsing
	content, err := os.ReadFile(config_file)
	handle_err(err)

	content2 := string(content)
	content3 := strings.Split(content2, "\n")

	// Node creation
	self_node_name = node_name
	self_node_id, _ = strconv.Atoi(node_name[len(node_name)-1:])
	print(self_node_id)

	total_nodes, _ := strconv.Atoi(content3[0])

	for i := 1; i <= total_nodes; i++ {
		node_info := strings.Split(content3[i], " ")

		// node_id, ip_addr: get additional fields
		node_id, _ := strconv.Atoi(node_info[0][len(node_info[0])-1:])

		ip_addr_net, err := net.LookupIP(node_info[1])
		ip_addr := ip_addr_net[0].String()
		handle_err(err)

		new_node := node{
			node_name:         node_info[0],
			node_id:           node_id,
			host_name:         node_info[1],
			ip_addr:           ip_addr,
			port_num:          node_info[2],
			connected_to_self: false,
		}

		node_info_map[node_info[0]] = new_node
	}

	// Transaction setup
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for node_name := range node_info_map {
		// https://stackoverflow.com/questions/42605337/cannot-assign-to-struct-field-in-a-map
		job_queue_at_node := job_queues[node_name]
		job_queue_at_node.cond = sync.NewCond(job_queues[node_name].mutex)

		job_queues[node_name] = job_queue_at_node
	}

	// Connections
	total_conns = (total_nodes - 1) * 2
	self_node := node_info_map[node_name]

	// Threading Begins
	// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
	wg.Add(2)
	// Servers
	go recieve_conn_reqs(self_node.port_num)

	// Clients
	go send_conn_reqs(self_node.node_name)

	wg.Wait()
}

/////// 1) CONNECTIONS ///////

// Iterate through all the nodes that arent ourselves and establish a connection as client
func send_conn_reqs(self_name string) {
	// wg.Add(2)

	// print("in send conn reqs")
	for name, info := range node_info_map {
		if name != self_name {
			host := info.host_name
			port := info.port_num
			go send_req(host, port, name)
		}
	}

	// wg.Wait()
}

// in a single thread
// sends a request to establish connection
func send_req(host string, port string, name string) {
	var conn net.Conn

	print("in send req\n")
	for conn == nil {
		ip := host + ":" + port
		conn, err := net.Dial("tcp", ip)
		print("established connection in send req\n")

		// Use preexisting thread to handle new connection

		if err != nil {
			continue
		} else {
			wait_for_connections(conn, name, false)
			break
		}
		print("after wait for connections...\n")
	}

	// if error, then redial else wait for connections and break

	// print("outside send conn for loop\n")

}

func recieve_conn_reqs(port string) {
	print("recieving conn reqs...\n")

	for i := 0; i < total_conns/2; i++ {
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
	// easiest thing to do: keep two connections between two nodes -> one for listening, other for writing
	// Increment current number of connections
	print("At beginning of wait for connections...\n")

	curr_conns.mutex.Lock()
	curr_conns.curr_conns += 1
	curr_conns.mutex.Unlock()

	for curr_conns.curr_conns < total_conns {
		time.Sleep(20 * time.Millisecond)
	}

	sec := 5
	print("Found all connections. Sleeping for + " + strconv.Itoa(sec) + "seconds...\n")

	// Sleep for a few seconds - make sure all the other nodes have established connections
	// don't worry about multicast
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
	print("in handle recieving transactions\n")

	for {
		incoming, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(incoming))

		new_message := str_to_message(incoming)
		incoming_message_id := new_message.Message_id
		incoming_node_id := new_message.Origin_id
		incoming_message_proposals := new_message.Proposals

		_, ok := message_info_map.message_info_map[incoming_message_id]

		// Check if the message is in the dictionary or not, add if it isnt
		if !ok {
			message_info_map.mutex.Lock()
			message_info_map.message_info_map[incoming_message_id] = new_message
			message_info_map.mutex.Unlock()
		}

		old_message := message_info_map.message_info_map[incoming_message_id]

		// ISIS algo
		if incoming_node_id == self_node_id { // If origin is ourselves
			old_message.Proposals = combine_arrs(old_message.Proposals, incoming_message_proposals)

			// if proposals_arr = full
			if len(old_message.Proposals) == total_nodes-1 {
				// Determine final priority
				final_pri := max_arr(old_message.Proposals)
				old_message.Final_priority = final_pri

				// Put on jobqueue for sending nodes & deliver ourselves!
				old_message.DeliveredIds = combine_arrs_int(old_message.DeliveredIds, []int{self_node_id})

				for node_name := range node_info_map {
					if node_name != self_node_name {
						// Put on jobqueue

					}
				}
				deliver_message(old_message)

				// Update message in dictionary (combine proposals)
				message_info_map.message_info_map[incoming_message_id] = old_message

				// signal

			} else {
				// Update message in dictionary (combine proposals)
				message_info_map.message_info_map[incoming_message_id] = old_message
			}
		} else { // If origin was another node
			// Add proposal to Proposal array
			curr_proposal.mutex.Lock()

			curr_proposal.curr_proposal += 1
			curr_proposal.mutex.Unlock()
			// Add to jobqueue to be sent back to the original
		}

		print(incoming_message_id)
		print("\n")
		print(incoming_node_id)
		print("\n")
	}
}

func handle_sending_transactions(conn net.Conn, node_name string) {
	// print("in handle sending transactions\n")
	// look into condition vars, sleep/wakeup on the condition variable
	// since we have one consumer
	job_queues[node_name].mutex.Lock()

	curr_job_queue := job_queues[node_name].job_queue

	for len(curr_job_queue) <= 0 {
		print("No more jobs to send at " + node_name)
		job_queues[node_name].cond.Wait()
	}
	curr_job := curr_job_queue[0] // message struct
	conn.Write([]byte(message_to_str(curr_job)))
	curr_job_queue = curr_job_queue[1:]

	job_queues[node_name].mutex.Unlock()
}

func deliver_message(m message) {

}

////// Helpers ///////

func combine_arrs_int(arr1 []int, arr2 []int) []int {
	slice := append(arr1, arr2...)

	keys := make(map[int]bool)
	list := []int{}

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

func remove_duplicates(slice []float64) []float64 {
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

////// PRIORITY QUEUE //////

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

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *heap_message, message_id string, priority int) {
	item.message_id = message_id
	item.priority = priority
	heap.Fix(pq, item.index)
}

/////// Extraneous //////

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
