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
	if number <= -10 {
		fmt.Println("range1")
	} else if number >= -10 && number < 0 {
		fmt.Println("range2")
	} else {
		fmt.Println("range3")
	}
}
