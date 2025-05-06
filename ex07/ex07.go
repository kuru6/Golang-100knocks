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
	if number == 0 {
		fmt.Println("入力された数字は0")
	} else {
		fmt.Println("入力された数字は0ではない")
	}

}
