package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("数字を2個入力してください")
	a := numberinput()
	b := numberinput()
	c := 1
	for c != 0 {
		c = a % b
		a, b = b, c
	}
	if b == 1 {
		fmt.Printf("互いに素である")
	} else {
		fmt.Printf("互いに素ではない")
	}

}

func numberinput() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	num, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("数字を入力してください")
	}
	return (num)
}
