package main

import (
	"fmt"
	"testing"
)

func TestRetrieveImg(t *testing.T) {
	result := retrieveImages("127539803930077")
	for _, x := range result {
		fmt.Println(x)
	}
}
