package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	numbers := make([]int, 3)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sep := strings.Split(scanner.Text(), " ")
	numbers[0], _ = strconv.Atoi(sep[0])
	numbers[1], _ = strconv.Atoi(sep[1])
	numbers[2], _ = strconv.Atoi(sep[2])
	small := 2147483647
	big := -2147483648
	for i := 0; i <= 2; i++ {
		if small >= numbers[i] {
			small = numbers[i]
		}
		if big <= numbers[i] {
			big = numbers[i]
		}
	}
	for j := 0; j <= 2; j++ {
		if numbers[j] != small && numbers[j] != big {
			fmt.Println(numbers[j])
		}
	}
}
