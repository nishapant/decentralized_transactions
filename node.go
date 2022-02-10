package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

func main() {
    content, err := os.ReadFile("config.txt")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(content))

    content2 := string(content)

    content3 := strings.Split(content2, "\n")

    fmt.Println(content3[0])
    nodes_num, i := strconv.Atoi(content3[0])
    // fmt.Println(string(content))

}