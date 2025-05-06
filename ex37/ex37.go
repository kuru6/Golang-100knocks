package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	numbers := []int{3, 7, 0, 8, 4, 1, 9, 6, 5, 2}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	number, _ := strconv.Atoi(text)
	youso := numbers[number]
	fmt.Println(numbers[youso])
}
