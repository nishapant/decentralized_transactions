package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
)

// https://www.calhoun.io/creating-random-strings-in-go/

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		print(text)

		message_id := strconv.Itoa(rand.Int())

		print(message_id)

		// add_transactions_to_queues()
	}
}

func add_transactions_to_queues(self_name string) {
	// read from stdin
	for {
		reader := bufio.NewReader(os.Stdin)
		curr_transaction, _ := reader.ReadString('\n')
		if curr_transaction == "" {
			continue
		}

		// print(curr_transaction)

		message_id := strconv.Itoa(rand.Int())

		print(message_id)

	}
}
