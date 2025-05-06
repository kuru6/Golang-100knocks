package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var result int
	scanner := bufio.NewScanner(os.Stdin)
	for i := 1; i <= 5; i++ {
		scanner.Scan()
		text := scanner.Text()
		number, _ := strconv.Atoi(text)
		result += number
	}
	fmt.Println(result)
}
