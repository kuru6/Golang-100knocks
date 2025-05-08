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
	for i := 1; number != 1; i++ {
		if number%2 == 0 {
			number = number / 2
		} else {
			number = number*3 + 1
		}
		fmt.Printf("%d回目：%d\n", i, number)
	}
}
