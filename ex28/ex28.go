package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	number, _ := strconv.Atoi(text)
	result := 1
	if number <= 0 {
		result = 1
	}
	for i := number; i >= 1; i-- {
		result *= i
	}
	fmt.Println(result)
}
