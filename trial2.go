package main

import (
	"bufio"
	"os"
)

// https://www.calhoun.io/creating-random-strings-in-go/

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		print(text)
	}
}
