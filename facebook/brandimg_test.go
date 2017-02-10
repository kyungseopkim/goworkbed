package main

import (
	"fmt"
	"testing"
)

func TestRetrieveImg(t *testing.T) {
	result := retrieveImages("179016535481")
	for _, x := range result {
		fmt.Println(x)
	}
}
