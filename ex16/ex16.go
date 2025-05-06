package main

import (
	"bufio"
	"os"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		if text == "0" {
			break
		}
	}
}
