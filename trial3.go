package main

import "encoding/json"

// https://www.calhoun.io/creating-random-strings-in-go/

// type message struct {
// 	Data           string    `json:"Data"`
// 	Origin_id      int       `json:"Origin_id"`
// 	Proposals      []float64 `json:"Proposals"`      // null or data
// 	Message_id     string    `json:"Message_id"`     // hash
// 	Final_priority float64   `json:"Final_priority"` // null when start
// }

type message struct {
	Data           string
	Origin_id      int
	Proposals      []float64 // null or data
	Message_id     string    // hash
	Final_priority float64   // null when start
}

func main() {
	m := message{Data: "hi", Proposals: []float64{1.2, 2.3}, Origin_id: 3, Message_id: "Awefawef", Final_priority: -1}
	print(message_to_str(m))
	mes_str := message_to_str(m)
	print("here: " + str_to_message(mes_str).Message_id)
}

func str_to_message(m_str string) message {
	var m message
	print("hi\n")
	print(m_str)
	print("\n")
	err := json.Unmarshal([]byte(m_str), &m)
	if err != nil {
		print("no")
	}

	return m
}

func message_to_str(m message) string {
	m_json, _ := json.Marshal(m)
	m_str := string(m_json) + "\n"

	return m_str
}
