package main

import "time"

// https://www.calhoun.io/creating-random-strings-in-go/

func main() {
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	text, _ := reader.ReadString('\n')
	// 	print(text)

	// 	message_id := strconv.Itoa(rand.Int())

	// 	print(message_id)

	// 	// add_transactions_to_queues()
	// }
	t1 := time.Now().Unix()
	time.Sleep(5 * time.Second)
	t2 := time.Now().Unix()

	print(t2 - t1)

	// final := subtractTime(t1, t2)
	// print(final.)
}

func subtractTime(time1, time2 time.Time) float64 {
	diff := time2.Sub(time1).Seconds()
	return diff
}

// func add_transactions_to_queues(self_name string) {
// 	// read from stdin
// 	for {
// 		reader := bufio.NewReader(os.Stdin)
// 		curr_transaction, _ := reader.ReadString('\n')
// 		if curr_transaction == "" {
// 			continue
// 		}

// 		// print(curr_transaction)

// 		message_id := strconv.Itoa(rand.Int())

// 		print(message_id)

// 	}
// }
