package main

import "fmt"

func main() {
	numbers := []int{}
	for i := 0; i <= 9; i++ {
		numbers = append(numbers, i)
	}
	fmt.Println(numbers)
}
