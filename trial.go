package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
func removeDuplicateValues(slice []float64) []float64 {
	keys := make(map[float64]bool)
	list := []float64{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	// mess := message{Data: "afeawef", DeliveredIds: []int{1, 2, 4}, OriginId: 1, Proposals: []float64{1.1, 2.2, 3.3}, Message_id: "awefawef", Final_priority: -1.0}
	// // print(mess)
	// m_json, _ := json.Marshal(mess)
	// print(m_json)
	// m_str := string(m_json)
	// print(m_str)
	// print(message_to_str(mess))
	// strings.Join(stringArray," ")
	arr := removeDuplicateValues(append([]float64{1.2, 2.2, 3.1}, []float64{2.2, 5.4, 6.2}...))
	print(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), " "), "[]"))
	// print(strings.Join(arr, " "))
}
