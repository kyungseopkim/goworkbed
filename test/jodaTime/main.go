package main

import (
	"fmt"

	"github.com/vjeantet/jodaTime"
)

func main() {
	date, err := jodaTime.Parse("yyyyMMdd", "20161130")
	if err != nil {
		panic(err)
	}

	fmt.Println(date)
}
