package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("西暦を入力してください:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	year, _ := strconv.Atoi(text)
	if year%4 == 0 && year%100 != 0 {
		fmt.Printf("%dはうるう年である", year)
	} else if year%400 == 0 {
		fmt.Printf("%dはうるう年である", year)
	} else {
		fmt.Printf("%dはうるう年ではない", year)
	}
}
