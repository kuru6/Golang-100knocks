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

	if number <= 0 {
		return
	}
	for i := 1; i <= number; i++ {
		fmt.Print("*")
		if i%5 == 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
