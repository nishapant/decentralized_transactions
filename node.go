package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type node struct {
	node_name string
	host_name string
	port_num  string
}

func main() {
	args := os.Args[1:]
	node_name := args[0]
	config_file := args[1]

	content, err := os.ReadFile(config_file)
	handle_err(err)

	content2 := string(content)
	content3 := strings.Split(content2, "\n")

	fmt.Println(content3[0])
	nodes_num, err := strconv.Atoi(content3[0])

	m := make(map[string]node)

	for i := 1; i <= nodes_num; i++ {
		node_info := strings.Split(content3[i], " ")
		
		new_node := node{
			node_name: node_info[0],
			host_name: node_info[1],
			port_num:  node_info[2],
		}

		m[node_info[0]] = new_node
	}

	fmt.Println(m)

	self_node := m[node_name]

	// thread 1
	serv_port := ":" + self_node.port_num
	ln, err := net.Listen("tcp", serv_port)
	handle_err(err)

	// closes the listener, after the function returns 
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		handle_err(err)
		print(conn)
		// go handle_receive(conn)
	}

	// handle any incoming connections

	// thread 2
	// for k, v := range m {
	// 	if k != node_name {
	// 		conn, err = net.Dial("tcp", strings.Join([]string{v.host_name, v.port_num}, ":"))
	// 	}
	// }
}

func handle_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
