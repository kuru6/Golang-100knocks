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
	num, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("0〜65535の整数値を入力してください")
	}
	fmt.Printf("%016b\n", num)
}
