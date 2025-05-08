package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	a := scantext()
	b := scantext()
	c := scantext()
	var D int
	D = b*b - 4*a*c
	if D == 0 {
		fmt.Println("重解")
	} else if D < 0 {
		fmt.Println("2つの虚数解")
	} else {
		fmt.Println("2つ実数解")
	}
}

func scantext() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	num, _ := strconv.Atoi(text)
	return (num)
}
