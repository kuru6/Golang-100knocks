package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	numbers := make([]int, 5)
	for i := 0; i <= 4; i++ {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		number, _ := strconv.Atoi(text)
		numbers[i] = number
	}
	fmt.Println(numbers)
}
