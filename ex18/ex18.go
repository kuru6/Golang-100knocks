package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	numbers := make([]int, 10)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	number, _ := strconv.Atoi(text)
	for i := 0; i <= 9; i++ {
		numbers[i] = number
	}
	fmt.Println(numbers)
}
