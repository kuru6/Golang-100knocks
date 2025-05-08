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
	if number >= 1 && number <= 9 {
		fmt.Printf("%dis a single figure.", number)
	} else {
		fmt.Printf("%d is not single figure.", number)
	}

}
