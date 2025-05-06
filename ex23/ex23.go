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
	if -5 <= number && number <= 10 {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
