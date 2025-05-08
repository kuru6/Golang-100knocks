package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	Text := [3]string{"0", "1", "2"}
	Numbers := [3]int{0, 1, 2}
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i <= 2; i++ {
		scanner.Scan()
		text := scanner.Text()
		Text[i] = text
		number, _ := strconv.Atoi(Text[i])
		Numbers[i] = number
	}
	if Numbers[0] <= Numbers[1] && Numbers[1] <= Numbers[2] {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
