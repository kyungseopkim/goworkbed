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

func ExampleInnoTime_Floor() {
	time := InnoTime{true, 12, 26}
	fmt.Println(time.Floor(30))
	fmt.Println(time.Floor(15))
	fmt.Println(time.Floor(10))
	fmt.Println(time.Floor(5))
	// Output: 12:00
	// 12:15
	// 12:20
	// 12:25
}

func ExampleInnoTime_EnumerateTo() {
	begin := InnoTime{true, 1, 10}
	end := InnoTime{true, 2, 35}

	data := begin.EnumerateTo(end, 30)
	fmt.Println(data)
	// Output: [01:00 01:30 02:00 02:30]
}
