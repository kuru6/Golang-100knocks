package main

import (
	"fmt"
)

func main() {
	numbers := [10]int{3, 7, 0, 8, 4, 1, 9, 6, 5, 2}

	for i := 0; i < 9; i++ {
		diff := numbers[i] - numbers[i+1]
		fmt.Println(diff)
	}
}
