package main

import "fmt"

// Move function moves
//
func Move(n, from, via, to int) {
    if n <= 0 {
      return
    }
    Move(n-1, from, to, via)
    fmt.Printf("Disk %d from %d to %d\n", n, from, to)
    Move(n-1, via, from, to)
}

// Hanoi behaves hanoi tower algorithm
//
func Hanoi(n int) {
  fmt.Println("Number of Disks:", n)
  Move(n, 1, 2, 3)
}

func main() {
  Hanoi(3)
}
