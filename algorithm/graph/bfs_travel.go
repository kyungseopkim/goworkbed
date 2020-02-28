package main

import (
	"fmt"
)

/*
data : [ (1, 2), (1, 3), (2, 7), (3, 4), (3, 6), (4, 5), (4, 6), (5, 6), (5, 7), (5, 9), (6, 8), (6, 9), (7, 8), (8, 9)]
*/

type tuple [2]int16

func buildGraph(data []tuple) map[int][]int {

	return nil
}

func main() {
	var data []tuple = []tuple{tuple{1, 2}, tuple{1, 3}, tuple{2, 7},
		tuple{3, 4}, tuple{3, 6}, tuple{4, 5},
		tuple{4, 6}, tuple{5, 6}, tuple{5, 7},
		tuple{5, 9}, tuple{6, 8}, tuple{6, 9},
		tuple{7, 8}, tuple{8, 9}}
	fmt.Println(data)

//	graph := make(map[int16][]int)

}
