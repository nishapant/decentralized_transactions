package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

type node struct {
    node_name string
    host_name string
    port_num int
  }

func main() {
    content, err := os.ReadFile("config.txt")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(content))

    content2 := string(content)

    content3 := strings.Split(content2, "\n")

    fmt.Println(content3[0])
    nodes_num, err := strconv.Atoi(content3[0])
    fmt.Println(string(nodes_num))

    m := make(map[string]node)

    for i := 1; i <= nodes_num; i++ {
        node_info := strings.Split(content3[i], " ")
        port, port_err := strconv.Atoi(node_info[2])
        if port_err != nil {
            log.Fatal(port_err)
        }
        new_node := node {
            node_name: node_info[0],
            host_name: node_info[1],
            port_num: port,
        }

        m[node_info[0]] = new_node
    }

    fmt.Println(m)

}