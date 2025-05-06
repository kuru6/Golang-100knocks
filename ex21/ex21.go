package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	if "5" < text && text <= "10" {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
