package main

import (
	"encoding/json"
	"fmt"
)

func ExampleInnoDate_tostring() {
	val := InnoDate{true, 2000, 11, 12}
	fmt.Println(val)
	//Output: 2000-11-12
}

func ExampleInnoDate_json_marshal() {
	val := InnoDate{true, 2010, 1, 21}
	dat, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dat))
	//Output: "2010-01-21"
}

func ExampleInnoDate_json_unmarshal() {
	val := InnoDate{true, 2010, 1, 21}
	dat, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
	}
	var result InnoDate
	err = json.Unmarshal(dat, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	//Output: 2010-01-21

}

func ExampleInnoTime() {
	now := InnoTime{true, 12, 20}
	fmt.Println(now.String())
	//Output: 12:20
}
