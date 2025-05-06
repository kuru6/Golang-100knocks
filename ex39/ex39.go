package main

import (
	"fmt"
)

func main() {
	numbers := []int{3, 7, 0, 8, 4, 1, 9, 6, 5, 2}
	for i := 0; i <= 8; i++ {
		fmt.Println(numbers[i] - numbers[i+1])
	}
}
