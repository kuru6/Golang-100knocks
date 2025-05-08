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
	a := scanner.Text()
	scanner.Scan()
	b := scanner.Text()
	anumber, _ := strconv.Atoi(a)
	bnumber, _ := strconv.Atoi(b)
	fmt.Println(anumber, bnumber)
	anumber, bnumber = bnumber, anumber
	fmt.Println(anumber, bnumber)
}
