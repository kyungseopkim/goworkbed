package main

import (
	"fmt"
	"testing"
)

func TestParseTime(t *testing.T) {
	time := parseTime("211")
	fmt.Println(time)
	time = parseTime("0233")
	fmt.Println(time)
}
