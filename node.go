package main

import (
	"bufio"
	"container/heap"
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
var node_id = 0
var total_nodes = 0
var node_info_map = make(map[string]node)

type node struct {
	node_name         string
	host_name         string
	port_num          string
	connected_to_self bool
	is_ready          bool // connected to all nodes in the network
}

// THREADING

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
	mutex          sync.Mutex
	data           string
	deliveredIds   []int // the process ids where the message is delivered
	originId       int
	proposals      []float64 // null or data
	message_id     string    // hash
	final_priority float64   // null when start
}

type heap_message struct {
	message_id string // hash
	index      int    // The index of the item in the heap.
	priority   int
	// The index is needed by update and is maintained by the heap.Interface methods.
}

type PriorityQueue []*heap_message

// var pq = PriorityQueue

type message_info_mutex struct {
	mutex            sync.Mutex
	message_info_map map[string]message
}

var message_info_map = message_info_mutex{
	message_info_map: make(map[string]message),
}

func main() {
	// time.Sleep(5 * time.Second)
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
	node_id := node_name[:len(node_name)-1]
	total_nodes, _ := strconv.Atoi(content3[0])

	for i := 1; i <= total_nodes; i++ {
		node_info := strings.Split(content3[i], " ")

		new_node := node{
			node_name:         node_info[0],
			host_name:         node_info[1],
			port_num:          node_info[2],
			connected_to_self: false,
		}

		node_info_map[node_info[0]] = new_node
	}

	// Transaction setup
	pq := make(PriorityQueue, 0)

	heap.Init(&pq)

	// Connections
	total_conns = (total_nodes - 1) * 2
	self_node := node_info_map[node_name]

	// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
	wg.Add(2)

	// print("going to recieve\n")
	// Server
	go recieve_conn_reqs(self_node.port_num)

	// print("going to send\n")
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
			go send_req(host, port)
		}
	}

	// wg.Wait()
}

// in a single thread
// sends a request to establish connection
func send_req(host string, port string) {
	var conn net.Conn

	// print("in send req\n")
	for conn == nil {
		ip := host + ":" + port
		conn, err := net.Dial("tcp", ip)
		// print("established connection in send req\n")

		// // Use preexisting thread to handle new connection
		// print("AFTER WAIT FOR CONNECTION...\n")
		// print(conn)
		// print("\n")

		if err != nil {
			continue
		} else {
			wait_for_connections(conn, false)
			break
		}
		// print("after wait for connections...\n")
	}

	// if error, then redial else wait for connections and break

	print("outside send conn for loop\n")

}

func recieve_conn_reqs(port string) {
	print("recieving conn reqs...\n")
	// wg.Add(2)

	for i := 0; i < total_conns/2; i++ {
		go recieve_req(port)
	}
	print("finished making threads for conn reqs...\n")
	// wg.Wait()
}

func recieve_req(port string) {
	print("recieving...\n")
	// Listen for incoming connections
	serv_port := ":" + port
	ln, err := net.Listen("tcp", serv_port)
	handle_err(err)

	// Accept connecitons 
	conn, err := ln.Accept()
	handle_err(err)

	// Close listener 
	ln.Close()

	wait_for_connections(conn, true)
}

func wait_for_connections(conn net.Conn, recieving bool) {
	// easiest thing to do: keep two connections between two nodes -> one for listening, other for writing
	// Increment current number of connections
	print("At beginning of wait for connections...\n")

	curr_conns.mutex.Lock()
	curr_conns.curr_conns += 1
	curr_conns.mutex.Unlock()

	// print("passed mutex...\n")
	print(curr_conns.curr_conns)
	// print(" curr cons\n")

	for curr_conns.curr_conns < total_conns {
		time.Sleep(20 * time.Millisecond)
	}

	sec := 5
	print("Found all connections. Sleeping for + " + strconv.Itoa(sec) + "seconds...\n")

	// Sleep for a few seconds to make sure all the other nodes have established connections
	// don't worry about multicast
	time.Sleep(5 * time.Second)

	// Move to handling transactions
	if recieving {
		handle_recieving_transactions(conn)
	} else {
		handle_sending_transactions(conn)
	}
}

////// 2) TRANSACTIONS  ///////

func handle_recieving_transactions(conn net.Conn) {
	print("in handle recieving transactions\n")
	print(conn)
	print("\n")

	for {
		incoming, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(incoming))
	}

}

func handle_sending_transactions(conn net.Conn) {
	print("in handle sending transactions\n")

	newmessage := "hi"
	// type message struct {
	// 	mutex          sync.Mutex
	// 	data           string
	// 	deliveredIds   []int // the process ids where the message is delivered
	// 	originId       int
	// 	proposals      []float64 // null or data
	// 	message_id     string    // hash
	// 	final_priority float64   // null when start
	// }

	newmessage = message {data: "awefewf->3", deliveredIds: [], originId: node_id}
	
	conn.Write([]byte(newmessage + "\n"))
}

////// Error Handling ///////

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
