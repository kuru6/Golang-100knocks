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
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	first_number, _ := strconv.Atoi(scanner.Text())
	min := first_number
	max := first_number
	for i := 1; i < n; i++ {
		scanner.Scan()
		num, _ := strconv.Atoi(scanner.Text())
		if num < min {
			min = num
		} else if num > max {
			max = num
		}
	}
	fmt.Println(min, max)
}
