package main

import "encoding/json"

// func main() {
// 	addr, err := net.LookupIP("sp22-cs425-g70-01.cs.illinois.edu")
// 	if err != nil {
// 		fmt.Println("Unknown host")
// 	} else {
// 		fmt.Println("IP address: ", addr)
// 	}
// }
type message struct {
	Data           string
	DeliveredIds   []int // the process ids where the message is delivered
	OriginId       int
	Proposals      []float64 // null or data
	Message_id     string    // hash
	Final_priority float64   // null when start
}

func message_to_str(m message) string {
	m_json, _ := json.Marshal(m)
	print(m_json)
	m_str := string(m_json)
	print(m_str)
	return m_str
}

func main() {
	mess := message{Data: "afeawef", DeliveredIds: []int{1, 2, 4}, OriginId: 1, Proposals: []float64{1.1, 2.2, 3.3}, Message_id: "awefawef", Final_priority: -1.0}
	// print(mess)
	m_json, _ := json.Marshal(mess)
	print(m_json)
	m_str := string(m_json)
	print(m_str)
	print(message_to_str(mess))
}
