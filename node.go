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

	print(node_name)

	content, err := os.ReadFile(config_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))

	content2 := string(content)

	content3 := strings.Split(content2, "\n")

	fmt.Println(content3[0])
	nodes_num, err := strconv.Atoi(content3[0])
	// fmt.Println(string(nodes_num))

	m := make(map[string]node)

	for i := 1; i <= nodes_num; i++ {
		node_info := strings.Split(content3[i], " ")
		// port, port_err := strconv.Atoi(node_info[2])
		// if port_err != nil {
		// 	log.Fatal(port_err)
		// }
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

	// handle any incoming connections

	// thread 2
	for k, v := range m {
		if k != node_name {
			conn, err = net.Dial("tcp", strings.Join([]string{v.host_name, v.port_num}, ":"))
		}
	}
}
