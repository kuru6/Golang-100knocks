package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	numbers := make([]int, 5)
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i <= 4; i++ {
		scanner.Scan()
		text := scanner.Text()
		num, _ := strconv.Atoi(text)
		numbers[i] = num
	}
	for j := 0; j <= 4; j++ {
		number := numbers[j]
		fmt.Print(j, "\t")
		for k := 1; k <= number; k++ {
			fmt.Print("*")
			if k%5 == 0 {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}

}
