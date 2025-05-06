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
	for i := 1; i <= 9; i++ {
		if i == number || i == number+1 {
			fmt.Print("")
		} else {
			fmt.Println(i)
		}
	}
}
