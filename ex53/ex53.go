package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("数値を入力してください")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	number, _ := strconv.Atoi(text)
	if number < 0 {
		fmt.Println("error")
	} else if number == 1 {
		fmt.Println("1")
	}
	for i := 2; i <= number; {
		if number%i == 0 {
			fmt.Print(i, "\t")
			number = number / i
		} else {
			i++
		}
	}
}
