package main

import (
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

func main() {
	time.Sleep(5 * time.Second)
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
	total_nodes, err := strconv.Atoi(content3[0])
	handle_err(err)

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

	// print("node info map:")
	// print(node_info_map)
	// print("\n")

	// Connections
	// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
	total_conns = (total_nodes - 1) * 2
	print("total_conns" + strconv.Itoa(total_conns))
	self_node := node_info_map[node_name]
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
			wait_for_connections(conn)
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
	serv_port := ":" + port
	ln, err := net.Listen("tcp", serv_port)
	handle_err(err)

	conn, err := ln.Accept()
	handle_err(err)

	wait_for_connections(conn)

	defer ln.Close()
}

func wait_for_connections(conn net.Conn) {
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
	// when done establishing connection, add a 1-3 second timeout to let all the nodes establish connections to each other
	// don't worry about multicast
	time.Sleep(5 * time.Second)

	// Move to handling transactions
	handle_transactions(conn)
}

////// 2) Transactions  ///////

func handle_transactions(conn net.Conn) {
	print("in handle transactions\n")
	print(conn)
	print()
}

////// Error Handling ///////

func handle_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
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
