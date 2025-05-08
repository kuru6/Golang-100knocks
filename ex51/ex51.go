package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("金額を入力してください: ")
	scanner.Scan()
	money, _ := strconv.Atoi(scanner.Text())

	onehundred := money / 100
	money %= 100

	ten := money / 10
	money %= 10

	one := money

	fmt.Printf("100円玉は%d枚, 10円玉は%d枚, 1円玉は%d枚\n", onehundred, ten, one)
}
