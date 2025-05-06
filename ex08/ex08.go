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
	if number < 0 {
		fmt.Println("数値は0より小さい")
	} else {
		fmt.Println("数値は0以上")
	}
}
