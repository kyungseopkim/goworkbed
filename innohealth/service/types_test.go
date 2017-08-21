package main

import (
	"encoding/json"
	"fmt"
)

//ExampleInnoWeekday is test
func ExampleInnoWeekday() {
	val := InnoWeekday(4)
	data, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
	// Output: "목"
}

