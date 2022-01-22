package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Rawdata string `json:"rawdata"`
}

func main() {

	jsonStr := `{"rawdata":"eyJyZXBsaWNhIjpbMSwyXX0="}`
	fmt.Printf("%#v\n", jsonStr)
	abc := []byte(jsonStr)
	fmt.Println("After byte", abc)
	var result User
	err := json.Unmarshal(abc, &result)
	fmt.Println("unmarshal the json obj", result)
	fmt.Printf("%#v", result)
	fmt.Println(err)
	var m []User

	m = append(m, result)
	fmt.Println(m)
	for v := range m {
		call(v)
	}
}
func call(v int) {
	fmt.Println(v)
	fmt.Printf("%#v\n", v)
}
