package main

import "encoding/json"

// https://www.calhoun.io/creating-random-strings-in-go/

type message struct {
	Data string
	// DeliveredIds   []int // the process ids where the message is delivered
	Origin_id      int
	Proposals      []float64 // null or data
	Message_id     string    // hash
	Final_priority float64   // null when start
}

func main() {
}

func str_to_message(m_str string) message {
	var m *message
	print("hi\n")
	print(m_str)
	print("\n")
	err := json.Unmarshal([]byte(m_str), m)
	handle_err(err)

	return *m
}
