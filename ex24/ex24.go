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
	if -10 <= number && number < 0 || 10 <= number {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
