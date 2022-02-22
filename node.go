package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
var total_conns = 0
var curr_conns = 0

// CONNECTIONS

type nodes_ready struct {
	mutex           sync.Mutex
	all_nodes_ready bool
}

var ready_flag = nodes_ready{all_nodes_ready: false}

// TRANSACTIONS

func main() {
	// Argument parsing
	if len(os.Args) < 3 {
		print("Incorrect number of Arguments!")
		os.Exit(1)
	}

	args := os.Args[1:]
	node_name := args[0]
	config_file := args[1]

	print("finished arg parsing")

	// File Parsing
	content, err := os.ReadFile(config_file)
	handle_err(err)

	content2 := string(content)
	content3 := strings.Split(content2, "\n")

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

	print(node_info_map)

	// 1) Connections
	// https://medium.com/@greenraccoon23/multi-thread-for-loops-easily-and-safely-in-go-a2e915302f8b
	total_conns = total_nodes - 1
	self_node := node_info_map[node_name]

	print("going to recieve")
	// Server
	recieve_conn_reqs(self_node.port_num)

	print("going to send")
	// Clients
	send_conn_reqs(self_node.node_name)

	// 2) Transactions

}

/////// 1) CONNECTIONS ///////

// Iterate through all the nodes that arent ourselves and establish a connection as client
func send_conn_reqs(self_name string) {
	wg.Add(total_conns)

	print("in send conn reqs")
	for name, info := range node_info_map {
		if name != self_name {
			host := info.host_name
			port := info.port_num
			go send_req(host, port)
		}
	}

	wg.Wait()
}

// in a single thread
// sends a request to establish connection
func send_req(host string, port string) {
	var conn net.Conn

	print("in send req")
	for conn == nil {
		ip := host + ":" + port
		conn, err := net.Dial("tcp", ip)

		curr_conns += 1

		// Use preexisting thread to handle new connection
		handle_connection(conn)

		if err != nil {
			continue
		}
		print(conn)
	}

}

func recieve_conn_reqs(port string) {
	serv_port := ":" + port
	ln, err := net.Listen("tcp", serv_port)
	handle_err(err)

	print("in recieve conn reqs")

	for curr_conns < total_conns {
		conn, err := ln.Accept()
		handle_err(err)

		curr_conns += 1

		// Start thread for connection
		go handle_connection(conn)
	}

	wg.Wait()

	// Close the listener
	defer ln.Close()
}

func handle_connection(conn net.Conn) {
	// Check if all connections are ready
	// ready_flag.mutex.Lock()
	// // use a mutex
	// for !ready_flag.all_nodes_ready {

	// }

	// ready_flag.mutex.Unlock()

}

////// 2) Transactions  ///////

////// Error Handling ///////

func handle_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
