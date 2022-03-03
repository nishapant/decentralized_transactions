package main

import "strconv"

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
	// random_hash()
	// m := message{Data: "hi", Proposals: []float64{1.2, 2.3}, Origin_id: 3, Message_id: "Awefawef", Final_priority: -1}
	// print(message_to_str(m))
	// mes_str := message_to_str(m)
	// print("here: " + str_to_message(mes_str).Message_id)

	account := info[1]
	print("info:", info, "\n")
	print("info2:", info[2], "\n")
	amount, _ := strconv.Atoi(info[2].)
}

// func str_to_message(m_str string) message {
// 	var m message
// 	print("hi\n")
// 	print(m_str)
// 	print("\n")
// 	err := json.Unmarshal([]byte(m_str), &m)
// 	if err != nil {
// 		print("no")
// 	}

// 	return m
// }

// func message_to_str(m message) string {
// 	m_json, _ := json.Marshal(m)
// 	m_str := string(m_json) + "\n"

// 	return m_str
// }

// func random_hash() string {
// 	num_arr := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
// 	hash := ""
// 	for i := 0; i <= 64; i++ {
// 		rand_int := rand.Intn(len(num_arr))
// 		hash += string(num_arr[rand_int])
// 	}

// 	print(hash)

// 	return hash
// }

// func deposit(info []string) {
// 	account := info[1]
// 	print("info:", info, "\n")
// 	print("info2:", info[2], "\n")
// 	amount, _ := strconv.Atoi(info[2])

// 	print("amount depositing: ", amount, "\n")

// 	bank.mutex.Lock()
// 	_, account_ok := bank.bank[account]
// 	if !account_ok {
// 		bank.bank[account] = 0
// 	}

// 	print("account: ", account, "\n")
// 	print("total before deposit: ", bank.bank[account], "\n")

// 	bank.bank[account] += amount

// 	print("total after: ", bank.bank[account], "\n")
// 	bank.mutex.Unlock()
// }
